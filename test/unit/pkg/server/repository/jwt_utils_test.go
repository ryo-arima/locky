package repository

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "test-secret"
	email := "test@example.com"

	// Test token generation
	tokenString, err := repository.GenerateJWT(email, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Test token validation
	claims, err := repository.ValidateJWT(tokenString, secret)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check claims
	assert.Equal(t, email, claims["email"])
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.True(t, time.Now().Unix() < int64(exp))
}

func TestValidateJWTWithInvalidToken(t *testing.T) {
	secret := "test-secret"

	// Test with an invalid token string
	_, err := repository.ValidateJWT("invalid-token", secret)
	assert.Error(t, err)

	// Test with a token signed with a different secret
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": "test@example.com",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte("different-secret"))
	assert.NoError(t, err)

	_, err = repository.ValidateJWT(token, secret)
	assert.Error(t, err)
	assert.Equal(t, "signature is invalid", err.Error())
}
