package repository

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// createMemberTestConfig creates a minimal BaseConfig for testing
func createMemberTestConfig() config.BaseConfig {
	return config.BaseConfig{
		DBConnection: nil, // No database connection needed for these tests
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Common: config.Common{},
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@test.com"},
					},
				},
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
					UserEmail:      "test@test.com",
					UserPassword:   "testpass",
				},
			},
			MySQL: config.MySQL{
				Host: "127.0.0.1",
				User: "root",
				Pass: "mysql",
				Port: "3306",
				Db:   "locky_test",
			},
		},
	}
}

// TestMemberRepository_NewMemberRepository tests repository creation without database dependency
func TestMemberRepository_NewMemberRepository(t *testing.T) {
	// Setup
	baseConfig := createMemberTestConfig()

	// Test
	memberRepo := repository.NewMemberRepository(baseConfig)

	// Assert using go-cmp with nil comparison
	if diff := cmp.Diff(false, memberRepo == nil); diff != "" {
		t.Errorf("Repository creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional verification - check if the repository has expected methods
	// This uses interface satisfaction check
	var _ repository.MemberRepository = memberRepo
}

// TestMemberRepository_InterfaceCompliance tests that the repository implements the expected interface
func TestMemberRepository_InterfaceCompliance(t *testing.T) {
	// Setup
	baseConfig := config.BaseConfig{
		DBConnection: nil,
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Common: config.Common{},
				Server: config.Server{
					Admin: config.Admin{
						Emails: []string{"admin@test.com"},
					},
				},
				Client: config.Client{
					ServerEndpoint: "http://localhost:8080",
				},
			},
		},
	}

	// Test
	memberRepo := repository.NewMemberRepository(baseConfig)

	// Assert - verify interface compliance
	var _ repository.MemberRepository = memberRepo

	// Test that the repository is not nil
	if diff := cmp.Diff(false, memberRepo == nil); diff != "" {
		t.Errorf("Repository should not be nil (-want +got):\n%s", diff)
	}
}

// TestMemberRepository_TableDriven tests repository creation with different configurations
func TestMemberRepository_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		dbConnection interface{}
		expectNil    bool
		expectError  bool
	}{
		{
			name:         "nil database connection",
			dbConnection: nil,
			expectNil:    false, // Repository should still be created
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			baseConfig := config.BaseConfig{
				DBConnection: nil,
				YamlConfig: config.YamlConfig{
					Application: config.Application{
						Common: config.Common{},
						Server: config.Server{
							Admin: config.Admin{
								Emails: []string{"admin@test.com"},
							},
						},
						Client: config.Client{
							ServerEndpoint: "http://localhost:8080",
						},
					},
				},
			}

			// Test
			memberRepo := repository.NewMemberRepository(baseConfig)

			// Assert
			isNil := memberRepo == nil
			if diff := cmp.Diff(tt.expectNil, isNil); diff != "" {
				t.Errorf("Repository nil check mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestMemberRepository_ConfigurationValidation tests that different configurations work
func TestMemberRepository_ConfigurationValidation(t *testing.T) {
	// Test different MySQL configurations
	configs := []config.MySQL{
		{
			Host: "localhost",
			User: "test",
			Pass: "test",
			Port: "3306",
			Db:   "test_db",
		},
		{
			Host: "127.0.0.1",
			User: "root",
			Pass: "mysql",
			Port: "3306",
			Db:   "locky_test",
		},
	}

	for i, mysqlConfig := range configs {
		t.Run(fmt.Sprintf("config_%d", i), func(t *testing.T) {
			baseConfig := config.BaseConfig{
				DBConnection: nil,
				YamlConfig: config.YamlConfig{
					Application: config.Application{
						Common: config.Common{},
						Server: config.Server{
							Admin: config.Admin{
								Emails: []string{"admin@test.com"},
							},
						},
						Client: config.Client{
							ServerEndpoint: "http://localhost:8080",
						},
					},
					MySQL: mysqlConfig,
				},
			}

			// Test
			memberRepo := repository.NewMemberRepository(baseConfig)

			// Assert
			if memberRepo == nil {
				t.Error("Expected MemberRepository to be created")
			}
		})
	}
}

// TestMemberRepository_EdgeCases tests edge cases for repository creation
func TestMemberRepository_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		baseConfig config.BaseConfig
	}{
		{
			name: "minimal config",
			baseConfig: config.BaseConfig{
				DBConnection: nil,
				YamlConfig:   config.YamlConfig{},
			},
		},
		{
			name: "empty admin emails",
			baseConfig: config.BaseConfig{
				DBConnection: nil,
				YamlConfig: config.YamlConfig{
					Application: config.Application{
						Server: config.Server{
							Admin: config.Admin{
								Emails: []string{},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			memberRepo := repository.NewMemberRepository(tt.baseConfig)

			// Assert - repository should still be created even with minimal config
			if memberRepo == nil {
				t.Error("Expected MemberRepository to be created even with minimal config")
			}

			// Verify interface compliance
			var _ repository.MemberRepository = memberRepo
		})
	}
}
