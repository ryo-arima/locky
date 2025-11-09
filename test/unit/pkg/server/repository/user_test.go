package repository

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// TestConvertModelToRequest tests the utility function without database access
func TestConvertModelToRequest(t *testing.T) {
	// Setup
	now := time.Now()
	userModel := model.Users{
		ID:        1,
		UUID:      "test-uuid",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Name:      "Test User",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Expected result
	expectedRequest := map[string]interface{}{
		"Email":    "test@example.com",
		"Password": "hashedpassword",
		"Name":     "Test User",
	}

	// Test
	userRequest := repository.ConvertModelToRequest(userModel)

	// Assert using go-cmp for better comparison
	actualRequest := map[string]interface{}{
		"Email":    userRequest.Email,
		"Password": userRequest.Password,
		"Name":     userRequest.Name,
	}

	if diff := cmp.Diff(expectedRequest, actualRequest); diff != "" {
		t.Errorf("ConvertModelToRequest() mismatch (-want +got):\n%s", diff)
	}

	// Additional specific field checks with detailed error messages
	if userRequest.Email != userModel.Email {
		t.Errorf("Email field mismatch: want %q, got %q", userModel.Email, userRequest.Email)
	}
	if userRequest.Password != userModel.Password {
		t.Errorf("Password field mismatch: want %q, got %q", userModel.Password, userRequest.Password)
	}
	if userRequest.Name != userModel.Name {
		t.Errorf("Name field mismatch: want %q, got %q", userModel.Name, userRequest.Name)
	}
}

// TestConvertUUIDToRequest tests the utility function without database access
func TestConvertUUIDToRequest(t *testing.T) {
	tests := []struct {
		name string
		uuid string
		want string
	}{
		{
			name: "basic uuid conversion",
			uuid: "test-uuid-123",
			want: "test-uuid-123",
		},
		{
			name: "empty uuid",
			uuid: "",
			want: "",
		},
		{
			name: "uuid with special characters",
			uuid: "test-uuid-456-special@#$",
			want: "test-uuid-456-special@#$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			userRequest := repository.ConvertUUIDToRequest(tt.uuid)

			// Assert using go-cmp
			if diff := cmp.Diff(tt.want, userRequest.Email); diff != "" {
				t.Errorf("ConvertUUIDToRequest().Email mismatch (-want +got):\n%s", diff)
			}

			// Additional check
			if userRequest.Email != tt.uuid {
				t.Errorf("Expected email %q, got %q", tt.uuid, userRequest.Email)
			}
		})
	}
}

// TestUserRepository_NewUserRepository tests repository creation without database dependency
func TestUserRepository_NewUserRepository(t *testing.T) {
	// Setup - Create a minimal BaseConfig without actual database connection
	baseConfig := config.BaseConfig{
		DBConnection: nil, // No database connection needed for constructor test
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

	// Test
	userRepo := repository.NewUserRepository(baseConfig)

	// Assert using go-cmp with nil comparison
	if userRepo == nil {
		t.Error("Expected UserRepository to be created, got nil")
	}

	// Additional verification - check if the repository has expected methods
	// This uses interface satisfaction check
	var _ repository.UserRepository = userRepo

	// Using cmp to verify the repository is properly initialized
	if diff := cmp.Diff(true, userRepo != nil); diff != "" {
		t.Errorf("Repository creation check mismatch (-want +got):\n%s", diff)
	}
}

// TestConvertModelToRequest_TableDriven demonstrates table-driven tests with go-cmp
func TestConvertModelToRequest_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		userModel model.Users
		want      map[string]interface{}
	}{
		{
			name: "complete user model",
			userModel: model.Users{
				ID:       1,
				UUID:     "uuid-123",
				Email:    "user@example.com",
				Password: "password123",
				Name:     "John Doe",
			},
			want: map[string]interface{}{
				"Email":    "user@example.com",
				"Password": "password123",
				"Name":     "John Doe",
			},
		},
		{
			name: "empty user model",
			userModel: model.Users{
				ID:       0,
				UUID:     "",
				Email:    "",
				Password: "",
				Name:     "",
			},
			want: map[string]interface{}{
				"Email":    "",
				"Password": "",
				"Name":     "",
			},
		},
		{
			name: "user with special characters",
			userModel: model.Users{
				ID:       2,
				UUID:     "uuid-456",
				Email:    "test+tag@example.com",
				Password: "p@ssw0rd!",
				Name:     "Jane O'Connor",
			},
			want: map[string]interface{}{
				"Email":    "test+tag@example.com",
				"Password": "p@ssw0rd!",
				"Name":     "Jane O'Connor",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			got := repository.ConvertModelToRequest(tt.userModel)

			// Create actual result map for comparison
			actual := map[string]interface{}{
				"Email":    got.Email,
				"Password": got.Password,
				"Name":     got.Name,
			}

			// Assert using go-cmp
			if diff := cmp.Diff(tt.want, actual); diff != "" {
				t.Errorf("ConvertModelToRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestUserRepository_UtilityFunctions tests various utility functions without database
func TestUserRepository_UtilityFunctions(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		input     interface{}
		expected  interface{}
	}{
		{
			name:      "convert basic model",
			operation: "model_to_request",
			input: model.Users{
				Email:    "test@example.com",
				Password: "password",
				Name:     "Test User",
			},
			expected: map[string]string{
				"Email":    "test@example.com",
				"Password": "password",
				"Name":     "Test User",
			},
		},
		{
			name:      "convert uuid to request",
			operation: "uuid_to_request",
			input:     "test-uuid-123",
			expected:  "test-uuid-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.operation {
			case "model_to_request":
				userModel := tt.input.(model.Users)
				result := repository.ConvertModelToRequest(userModel)
				expectedMap := tt.expected.(map[string]string)

				if diff := cmp.Diff(expectedMap["Email"], result.Email); diff != "" {
					t.Errorf("Email mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(expectedMap["Password"], result.Password); diff != "" {
					t.Errorf("Password mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(expectedMap["Name"], result.Name); diff != "" {
					t.Errorf("Name mismatch (-want +got):\n%s", diff)
				}
			case "uuid_to_request":
				uuid := tt.input.(string)
				result := repository.ConvertUUIDToRequest(uuid)
				expected := tt.expected.(string)

				if diff := cmp.Diff(expected, result.Email); diff != "" {
					t.Errorf("UUID conversion mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
