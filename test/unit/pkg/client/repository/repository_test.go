package client

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

	repo := repository.NewCommon(cfg)
	assert.NotNil(t, repo)
}

func TestNewGroup(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewGroupInternal(cfg)
	assert.NotNil(t, repo)
}

func TestNewMember(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewMemberInternal(cfg)
	assert.NotNil(t, repo)
}

func TestNewRole(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewRoleInternal(cfg)
	assert.NotNil(t, repo)
}

func TestNewUser(t *testing.T) {
	cfg := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	repo := repository.NewUserInternal(cfg)
	assert.NotNil(t, repo)
}

func TestLoginRequest_Structure(t *testing.T) {
	loginReq := request.Login{
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
	var _ repository.Common = repository.NewCommon(cfg)
	var _ repository.GroupInternal = repository.NewGroupInternal(cfg)
	var _ repository.MemberInternal = repository.NewMemberInternal(cfg)
	var _ repository.RoleInternal = repository.NewRoleInternal(cfg)
	var _ repository.UserInternal = repository.NewUserInternal(cfg)
}
