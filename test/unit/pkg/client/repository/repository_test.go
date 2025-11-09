package repository_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/stretchr/testify/assert"
)

// Mock server setup is complex, so we test repository creation and basic structure
func TestNewCommonRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					UserEmail:      "test@example.com",
					UserPassword:   "password",
				},
			},
		},
	}

	repo := repository.NewCommonRepository(cfg)
	assert.NotNil(t, repo)
}

func TestNewGroupRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewGroupRepository(cfg)
	assert.NotNil(t, repo)
}

func TestNewMemberRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewMemberRepository(cfg)
	assert.NotNil(t, repo)
}

func TestNewRoleRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewRoleRepository(cfg)
	assert.NotNil(t, repo)
}

func TestNewUserRepository(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewUserRepository(cfg)
	assert.NotNil(t, repo)
}

func TestLoginRequest_Structure(t *testing.T) {
	loginReq := request.LoginRequest{
		Email:    "user@example.com",
		Password: "securepassword123",
	}

	assert.Equal(t, "user@example.com", loginReq.Email)
	assert.Equal(t, "securepassword123", loginReq.Password)
}

// Note: Full integration tests for Login, RefreshToken, etc. are in test/e2e
// Unit tests verify repository creation and basic structure
func TestRepositoryInterfaces(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	// Verify all repositories implement their interfaces
	var _ repository.CommonRepository = repository.NewCommonRepository(cfg)
	var _ repository.GroupRepository = repository.NewGroupRepository(cfg)
	var _ repository.MemberRepository = repository.NewMemberRepository(cfg)
	var _ repository.RoleRepository = repository.NewRoleRepository(cfg)
	var _ repository.UserRepository = repository.NewUserRepository(cfg)
}
