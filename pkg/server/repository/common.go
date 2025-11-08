package repository

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/mail"
	"golang.org/x/crypto/bcrypt"
)

type CommonRepository interface {
	GetBaseConfig() config.BaseConfig
	GenerateJWTToken(claims model.JWTClaims) (string, error)
	ValidateJWTToken(tokenString string) (*model.JWTClaims, error)
	ParseTokenUnverified(tokenString string) (*model.JWTClaims, error)
	IsTokenInvalidated(ctx context.Context, jti string) (bool, error)
	InvalidateToken(ctx context.Context, tokenString string) error
	GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error)
	GenerateJWTSecret() (string, error)
	ValidateJWTSecretStrength(secret string) error
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
	ValidatePasswordStrength(password string) error
	DeleteTokenCache(token string) // Added: Cache deletion on logout
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error
}

type commonRepository struct {
	BaseConfig  config.BaseConfig
	RedisClient *redis.Client
	MailSender  *mail.Sender
}

func (commonRepository *commonRepository) GetBaseConfig() config.BaseConfig {
	return commonRepository.BaseConfig
}

// GetJWTSecret returns the JWT secret key from config, environment variable, or a default value
func (cr *commonRepository) getJWTSecret() string {
	// First try environment variable
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		return envSecret
	}

	// Then try config file
	if cr.BaseConfig.YamlConfig.Application.Server.JWTSecret != "" {
		return cr.BaseConfig.YamlConfig.Application.Server.JWTSecret
	}

	// Finally, use default (should be changed in production)
	return "your-256-bit-secret-key-change-this-in-production"
}

// GenerateJWTToken creates a JWT token with the given claims
func (cr *commonRepository) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	// Create header
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(headerJSON)

	// Create payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(payloadJSON)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := cr.createSignature(message)

	// Combine all parts
	token := message + "." + signature
	return token, nil
}

// ValidateJWTToken validates and parses a JWT token
func (cr *commonRepository) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	// 1. Try cache first
	if cr.RedisClient != nil {
		if cached, err := cr.getCachedTokenClaims(tokenString); err == nil && cached != nil {
			// Ensure not expired
			if cached.ExpiresAt >= time.Now().Unix() {
				return cached, nil
			}
		}
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Verify signature
	message := parts[0] + "." + parts[1]
	expectedSignature := cr.createSignature(message)
	if parts[2] != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	// Decode payload
	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims model.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		return nil, errors.New("token expired")
	}

	// Cache claims (TTL = min(30m, remaining lifetime))
	if cr.RedisClient != nil {
		_ = cr.cacheTokenClaims(tokenString, &claims)
	}

	return &claims, nil
}

// GenerateTokenPair creates both access and refresh tokens
func (cr *commonRepository) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	now := time.Now()
	accessTokenExpiry := now.Add(24 * time.Hour).Unix()      // 24 hours
	refreshTokenExpiry := now.Add(7 * 24 * time.Hour).Unix() // 7 days

	// Generate a unique ID for the token
	jti := uuid.New().String()

	// Create access token claims
	accessClaims := model.JWTClaims{
		Jti:       jti,
		UserID:    userID,
		UUID:      userUUID,
		Email:     email,
		Name:      name,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: accessTokenExpiry,
	}

	// Generate access token
	accessToken, err := cr.GenerateJWTToken(accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create refresh token claims (longer expiry, different JTI)
	refreshJti := uuid.New().String()
	refreshClaims := model.JWTClaims{
		Jti:       refreshJti,
		UserID:    userID,
		UUID:      userUUID,
		Email:     email,
		Name:      name,
		Role:      role,
		IssuedAt:  now.Unix(),
		ExpiresAt: refreshTokenExpiry,
	}

	// Generate refresh token
	refreshToken, err := cr.GenerateJWTToken(refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessTokenExpiry - now.Unix(),
	}, nil
}

// createSignature creates HMAC-SHA256 signature for JWT
func (cr *commonRepository) createSignature(message string) string {
	secret := cr.getJWTSecret()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil))
	return signature
}

// HashPassword hashes a password using bcrypt
func (cr *commonRepository) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// VerifyPassword verifies a password against its hash
func (cr *commonRepository) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordStrength validates password strength
func (cr *commonRepository) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// GenerateJWTSecret generates a secure random JWT secret
func (cr *commonRepository) GenerateJWTSecret() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 for use as a string
	secret := base64.URLEncoding.EncodeToString(bytes)
	return secret, nil
}

// ValidateJWTSecretStrength checks if JWT secret meets minimum security requirements
func (cr *commonRepository) ValidateJWTSecretStrength(secret string) error {
	if len(secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// Check for common weak secrets
	weakSecrets := []string{
		"secret",
		"jwt-secret",
		"your-secret-key",
		"change-me",
		"your-256-bit-secret-key-change-this-in-production",
	}

	for _, weak := range weakSecrets {
		if secret == weak {
			return fmt.Errorf("JWT secret is too weak: %s", secret)
		}
	}

	return nil
}

// InvalidateToken adds a token's JTI to the Redis denylist
func (cr *commonRepository) InvalidateToken(ctx context.Context, tokenString string) error {
	claims, err := cr.ValidateJWTToken(tokenString)
	if err != nil {
		// If token is already expired or invalid, we don't need to do anything.
		// We can consider it "successfully" invalidated.
		if strings.Contains(err.Error(), "token expired") {
			return nil
		}
		return fmt.Errorf("error validating token before invalidation: %w", err)
	}

	// Calculate remaining time until expiration
	now := time.Now()
	expiresAt := time.Unix(claims.ExpiresAt, 0)
	if now.After(expiresAt) {
		// Already expired
		return nil
	}
	ttl := expiresAt.Sub(now)

	// Add to denylist in Redis
	err = cr.RedisClient.Set(ctx, claims.Jti, "invalidated", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to denylist: %w", err)
	}

	return nil
}

// ParseTokenUnverified decodes the claims from a token without verifying its signature.
// This is used to get the JTI for denylist checking before full validation.
func (cr *commonRepository) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payloadBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims model.JWTClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	return &claims, nil
}

// IsTokenInvalidated checks if a token's JTI exists in the Redis denylist.
func (cr *commonRepository) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	result, err := cr.RedisClient.Exists(ctx, jti).Result()
	if err != nil {
		return true, fmt.Errorf("error checking token in redis: %w", err)
	}
	return result == 1, nil
}

// helper: cache key builder
func (cr *commonRepository) tokenCacheKey(token string) string {
	return "auth:token:" + token
}

// helper: store token claims in redis with 30m max TTL
func (cr *commonRepository) cacheTokenClaims(token string, claims *model.JWTClaims) error {
	if cr.RedisClient == nil || claims == nil {
		return nil
	}
	data, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	remaining := time.Until(time.Unix(claims.ExpiresAt, 0))
	if remaining <= 0 {
		return nil
	}
	maxTTL := 30 * time.Minute
	if remaining < maxTTL {
		maxTTL = remaining
	}
	return cr.RedisClient.Set(context.Background(), cr.tokenCacheKey(token), string(data), maxTTL).Err()
}

// helper: get token claims from redis cache
func (cr *commonRepository) getCachedTokenClaims(token string) (*model.JWTClaims, error) {
	if cr.RedisClient == nil {
		return nil, errors.New("redis client nil")
	}
	val, err := cr.RedisClient.Get(context.Background(), cr.tokenCacheKey(token)).Result()
	if err != nil {
		return nil, err
	}
	var claims model.JWTClaims
	if err := json.Unmarshal([]byte(val), &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

// DeleteTokenCache removes cached claims for the given raw token string (if present)
func (cr *commonRepository) DeleteTokenCache(token string) {
	if cr.RedisClient == nil || token == "" {
		return
	}
	_ = cr.RedisClient.Del(context.Background(), cr.tokenCacheKey(token)).Err()
}

// SendEmail sends an email using the configured mail sender
func (cr *commonRepository) SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error {
	if cr.MailSender == nil {
		return errors.New("mail sender not configured")
	}

	msg := mail.Message{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  isHTML,
	}

	return cr.MailSender.Send(msg)
}

// SendWelcomeEmail sends a welcome email to a new user
func (cr *commonRepository) SendWelcomeEmail(ctx context.Context, to, name string) error {
	if cr.MailSender == nil {
		return errors.New("mail sender not configured")
	}
	return cr.MailSender.SendWelcomeEmail(to, name)
}

// SendPasswordResetEmail sends a password reset email to a user
func (cr *commonRepository) SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error {
	if cr.MailSender == nil {
		return errors.New("mail sender not configured")
	}
	return cr.MailSender.SendPasswordResetEmail(to, name, resetURL)
}

func NewCommonRepository(baseConfig config.BaseConfig, redisClient *redis.Client) CommonRepository {
	// Initialize mail sender from config
	var mailSender *mail.Sender
	if baseConfig.YamlConfig.Application.Mail.Host != "" {
		mailConfig := mail.Config{
			Host:     baseConfig.YamlConfig.Application.Mail.Host,
			Port:     baseConfig.YamlConfig.Application.Mail.Port,
			Username: baseConfig.YamlConfig.Application.Mail.Username,
			Password: baseConfig.YamlConfig.Application.Mail.Password,
			From:     baseConfig.YamlConfig.Application.Mail.From,
			UseTLS:   baseConfig.YamlConfig.Application.Mail.UseTLS,
		}
		mailSender = mail.NewSender(mailConfig)
	}

	return &commonRepository{
		BaseConfig:  baseConfig,
		RedisClient: redisClient,
		MailSender:  mailSender,
	}
}
