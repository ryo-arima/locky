package client

import (
	"testing"

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

	repo := share.NewCommonRepository(cfg)
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

	repo := share.NewGroup(cfg)
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

	repo := share.NewMember(cfg)
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

	repo := share.NewRole(cfg)
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

	repo := share.NewUser(cfg)
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
	var _ share.CommonRepository = share.NewCommonRepository(cfg)
	var _ share.GroupRepository = share.NewGroup(cfg)
	var _ share.MemberRepository = share.NewMember(cfg)
	var _ share.RoleRepository = share.NewRole(cfg)
	var _ share.UserRepository = share.NewUser(cfg)
}
