package usecase_test

import (
	"context"
	"testing"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommonUsecase(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	assert.NotNil(t, uc)
	assert.Equal(t, cfg, uc.GetBaseConfig())
}

func TestNewGroupUsecase(t *testing.T) {
	cfg := config.BaseConfig{}
	groupRepo := repository.NewGroupRepository(cfg)
	commonRepo := repository.NewCommonRepository(cfg, nil)

	uc := usecase.NewGroupUsecase(groupRepo, commonRepo)
	assert.NotNil(t, uc)
}

func TestNewMemberUsecase(t *testing.T) {
	cfg := config.BaseConfig{}
	memberRepo := repository.NewMemberRepository(cfg)
	commonRepo := repository.NewCommonRepository(cfg, nil)

	uc := usecase.NewMemberUsecase(memberRepo, commonRepo)
	assert.NotNil(t, uc)
}

func TestNewRoleUsecase(t *testing.T) {
	cfg := config.BaseConfig{}
	roleRepo := repository.NewRoleRepository(cfg)
	commonRepo := repository.NewCommonRepository(cfg, nil)

	uc := usecase.NewRoleUsecase(roleRepo, commonRepo)
	assert.NotNil(t, uc)
}

func TestNewUserUsecase(t *testing.T) {
	cfg := config.BaseConfig{}
	userRepo := repository.NewUserRepository(cfg)
	commonRepo := repository.NewCommonRepository(cfg, nil)

	uc := usecase.NewUserUsecase(userRepo, commonRepo)
	assert.NotNil(t, uc)
}

func TestCommonUsecase_HashPassword(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	password := "SecurePassword123!"
	hashed, err := uc.HashPassword(password)

	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestCommonUsecase_VerifyPassword(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	password := "SecurePassword123!"
	hashed, err := uc.HashPassword(password)
	require.NoError(t, err)

	// Correct password
	err = uc.VerifyPassword(hashed, password)
	assert.NoError(t, err)

	// Incorrect password
	err = uc.VerifyPassword(hashed, "WrongPassword")
	assert.Error(t, err)
}

func TestCommonUsecase_ValidatePasswordStrength(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid strong password",
			password: "StrongPass123!",
			wantErr:  false,
		},
		{
			name:     "Too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "No uppercase",
			password: "weakpass123!",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommonUsecase_GenerateJWTSecret(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	secret, err := uc.GenerateJWTSecret()
	require.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.GreaterOrEqual(t, len(secret), 32)
}

func TestCommonUsecase_ValidateJWTSecretStrength(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid secret",
			secret:  "this-is-a-very-long-and-secure-secret-key-for-jwt",
			wantErr: false,
		},
		{
			name:    "Too short",
			secret:  "short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateJWTSecretStrength(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommonUsecase_GenerateTokenPair(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	tokens, err := uc.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")

	require.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestCommonUsecase_ValidateJWTToken(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	// Generate a token
	tokens, err := uc.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")
	require.NoError(t, err)

	// Validate the token
	claims, err := uc.ValidateJWTToken(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
}

func TestCommonUsecase_ParseTokenUnverified(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key-for-jwt-that-is-long-enough",
				},
			},
		},
	}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	// Generate a token
	tokens, err := uc.GenerateTokenPair(1, "user-uuid-123", "test@example.com", "Test User", "user")
	require.NoError(t, err)

	// Parse without verification
	claims, err := uc.ParseTokenUnverified(tokens.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", claims.Email)
}

func TestCommonUsecase_IsTokenInvalidated(t *testing.T) {
	cfg := config.BaseConfig{}
	repo := repository.NewCommonRepository(cfg, nil)
	uc := usecase.NewCommonUsecase(repo)

	// Without Redis, this should return false (token is valid)
	invalidated, err := uc.IsTokenInvalidated(context.Background(), "test-jti")
	assert.NoError(t, err)
	assert.False(t, invalidated)
}
