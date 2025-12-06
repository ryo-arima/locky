package usecase

import (
	"context"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type Common interface {
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
	DeleteTokenCache(token string)
	SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error
	SendWelcomeEmail(ctx context.Context, to, name string) error
	SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error
}

type common struct {
	commonRepo repository.Common
}

func NewCommon(commonRepo repository.Common) Common {
	return &common{
		commonRepo: commonRepo,
	}
}

func (uc *common) GetBaseConfig() config.BaseConfig {
	return uc.commonRepo.GetBaseConfig()
}

func (uc *common) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	return uc.commonRepo.GenerateJWTToken(claims)
}

func (uc *common) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ValidateJWTToken(tokenString)
}

func (uc *common) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ParseTokenUnverified(tokenString)
}

func (uc *common) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	return uc.commonRepo.IsTokenInvalidated(ctx, jti)
}

func (uc *common) InvalidateToken(ctx context.Context, tokenString string) error {
	return uc.commonRepo.InvalidateToken(ctx, tokenString)
}

func (uc *common) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	return uc.commonRepo.GenerateTokenPair(userID, userUUID, email, name, role)
}

func (uc *common) GenerateJWTSecret() (string, error) {
	return uc.commonRepo.GenerateJWTSecret()
}

func (uc *common) ValidateJWTSecretStrength(secret string) error {
	return uc.commonRepo.ValidateJWTSecretStrength(secret)
}

func (uc *common) HashPassword(password string) (string, error) {
	return uc.commonRepo.HashPassword(password)
}

func (uc *common) VerifyPassword(hashedPassword, password string) error {
	return uc.commonRepo.VerifyPassword(hashedPassword, password)
}

func (uc *common) ValidatePasswordStrength(password string) error {
	return uc.commonRepo.ValidatePasswordStrength(password)
}

func (uc *common) DeleteTokenCache(token string) {
	uc.commonRepo.DeleteTokenCache(token)
}

func (uc *common) SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error {
	return uc.commonRepo.SendEmail(ctx, to, subject, body, isHTML)
}

func (uc *common) SendWelcomeEmail(ctx context.Context, to, name string) error {
	return uc.commonRepo.SendWelcomeEmail(ctx, to, name)
}

func (uc *common) SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error {
	return uc.commonRepo.SendPasswordResetEmail(ctx, to, name, resetURL)
}
