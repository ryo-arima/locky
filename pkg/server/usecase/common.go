package usecase

import (
	"context"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonUsecase interface {
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
}

type commonUsecase struct {
	commonRepo repository.CommonRepository
}

func NewCommonUsecase(commonRepo repository.CommonRepository) CommonUsecase {
	return &commonUsecase{
		commonRepo: commonRepo,
	}
}

func (uc *commonUsecase) GetBaseConfig() config.BaseConfig {
	return uc.commonRepo.GetBaseConfig()
}

func (uc *commonUsecase) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	return uc.commonRepo.GenerateJWTToken(claims)
}

func (uc *commonUsecase) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ValidateJWTToken(tokenString)
}

func (uc *commonUsecase) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	return uc.commonRepo.ParseTokenUnverified(tokenString)
}

func (uc *commonUsecase) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	return uc.commonRepo.IsTokenInvalidated(ctx, jti)
}

func (uc *commonUsecase) InvalidateToken(ctx context.Context, tokenString string) error {
	return uc.commonRepo.InvalidateToken(ctx, tokenString)
}

func (uc *commonUsecase) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	return uc.commonRepo.GenerateTokenPair(userID, userUUID, email, name, role)
}

func (uc *commonUsecase) GenerateJWTSecret() (string, error) {
	return uc.commonRepo.GenerateJWTSecret()
}

func (uc *commonUsecase) ValidateJWTSecretStrength(secret string) error {
	return uc.commonRepo.ValidateJWTSecretStrength(secret)
}

func (uc *commonUsecase) HashPassword(password string) (string, error) {
	return uc.commonRepo.HashPassword(password)
}

func (uc *commonUsecase) VerifyPassword(hashedPassword, password string) error {
	return uc.commonRepo.VerifyPassword(hashedPassword, password)
}

func (uc *commonUsecase) ValidatePasswordStrength(password string) error {
	return uc.commonRepo.ValidatePasswordStrength(password)
}

func (uc *commonUsecase) DeleteTokenCache(token string) {
	uc.commonRepo.DeleteTokenCache(token)
}
