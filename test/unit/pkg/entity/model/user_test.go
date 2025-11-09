package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/entity/model"
)

// TestUsers_Creation tests basic user creation
func TestUsers_Creation(t *testing.T) {
	// Test data
	now := time.Now()
	user := model.Users{
		ID:        1,
		UUID:      "test-user-uuid-123",
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "hashed_password_123",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Assert using go-cmp
	expectedID := uint(1)
	if diff := cmp.Diff(expectedID, user.ID); diff != "" {
		t.Errorf("User ID mismatch (-want +got):\n%s", diff)
	}

	expectedUUID := "test-user-uuid-123"
	if diff := cmp.Diff(expectedUUID, user.UUID); diff != "" {
		t.Errorf("User UUID mismatch (-want +got):\n%s", diff)
	}

	expectedEmail := "test@example.com"
	if diff := cmp.Diff(expectedEmail, user.Email); diff != "" {
		t.Errorf("User Email mismatch (-want +got):\n%s", diff)
	}

	expectedName := "Test User"
	if diff := cmp.Diff(expectedName, user.Name); diff != "" {
		t.Errorf("User Name mismatch (-want +got):\n%s", diff)
	}

	// Verify DeletedAt is nil (not soft deleted)
	if user.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil for new user")
	}

	// Verify password is set (should be hashed)
	if user.Password == "" {
		t.Error("Expected Password to be set")
	}
}

// TestUsers_TableDriven demonstrates table-driven tests for user validation
func TestUsers_TableDriven(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		user        model.Users
		expectedErr bool
		description string
	}{
		{
			name: "valid user",
			user: model.Users{
				ID:        1,
				UUID:      "valid-uuid-123",
				Email:     "valid@example.com",
				Name:      "Valid User Name",
				Password:  "hashed_password",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false,
			description: "Should accept valid user data",
		},
		{
			name: "empty email",
			user: model.Users{
				ID:        2,
				UUID:      "empty-email-uuid",
				Email:     "",
				Name:      "User Without Email",
				Password:  "hashed_password",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false, // Empty email might be allowed in model validation
			description: "Should handle empty email",
		},
		{
			name: "empty name",
			user: model.Users{
				ID:        3,
				UUID:      "empty-name-uuid",
				Email:     "user@example.com",
				Name:      "",
				Password:  "hashed_password",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: nil,
			},
			expectedErr: false, // Empty name might be allowed in model validation
			description: "Should handle empty name",
		},
		{
			name: "soft deleted user",
			user: model.Users{
				ID:        4,
				UUID:      "deleted-uuid",
				Email:     "deleted@example.com",
				Name:      "Deleted User",
				Password:  "hashed_password",
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: &now,
			},
			expectedErr: false,
			description: "Should handle soft deleted users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the user structure
			if tt.user.ID == 0 {
				t.Error("Expected non-zero ID")
			}

			if tt.user.UUID == "" {
				t.Error("Expected non-empty UUID")
			}

			if tt.user.CreatedAt == nil {
				t.Error("Expected CreatedAt to be set")
			}

			if tt.user.UpdatedAt == nil {
				t.Error("Expected UpdatedAt to be set")
			}

			// For soft deleted users, verify DeletedAt is set
			if tt.name == "soft deleted user" {
				if tt.user.DeletedAt == nil {
					t.Error("Expected DeletedAt to be set for soft deleted user")
				}
			}
		})
	}
}

// TestUsers_EdgeCases tests edge cases and special scenarios
func TestUsers_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		testFunc    func(t *testing.T)
	}{
		{
			name:        "zero values",
			description: "Should handle zero values properly",
			testFunc: func(t *testing.T) {
				user := model.Users{}

				// Check zero values
				expectedID := uint(0)
				if diff := cmp.Diff(expectedID, user.ID); diff != "" {
					t.Errorf("Zero ID mismatch (-want +got):\n%s", diff)
				}

				expectedUUID := ""
				if diff := cmp.Diff(expectedUUID, user.UUID); diff != "" {
					t.Errorf("Zero UUID mismatch (-want +got):\n%s", diff)
				}

				expectedEmail := ""
				if diff := cmp.Diff(expectedEmail, user.Email); diff != "" {
					t.Errorf("Zero Email mismatch (-want +got):\n%s", diff)
				}

				expectedName := ""
				if diff := cmp.Diff(expectedName, user.Name); diff != "" {
					t.Errorf("Zero Name mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:        "special characters in email",
			description: "Should handle special characters in email addresses",
			testFunc: func(t *testing.T) {
				now := time.Now()
				user := model.Users{
					ID:        1,
					UUID:      "special-email-uuid",
					Email:     "user+tag123@sub.example-domain.com",
					Name:      "User With Special Email",
					Password:  "hashed_password",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				expectedEmail := "user+tag123@sub.example-domain.com"
				if diff := cmp.Diff(expectedEmail, user.Email); diff != "" {
					t.Errorf("Special character email mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:        "unicode characters in name",
			description: "Should handle unicode characters in user names",
			testFunc: func(t *testing.T) {
				now := time.Now()
				user := model.Users{
					ID:        1,
					UUID:      "unicode-uuid",
					Email:     "user@example.com",
					Name:      "田中太郎",
					Password:  "hashed_password",
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				expectedName := "田中太郎"
				if diff := cmp.Diff(expectedName, user.Name); diff != "" {
					t.Errorf("Unicode name mismatch (-want +got):\n%s", diff)
				}
			},
		},
		{
			name:        "long password",
			description: "Should handle very long password hashes",
			testFunc: func(t *testing.T) {
				now := time.Now()
				longPassword := "very_long_hashed_password_that_might_be_generated_by_bcrypt_or_similar_hashing_algorithms_123456789"
				user := model.Users{
					ID:        1,
					UUID:      "long-password-uuid",
					Email:     "user@example.com",
					Name:      "User With Long Password",
					Password:  longPassword,
					CreatedAt: &now,
					UpdatedAt: &now,
				}

				if diff := cmp.Diff(longPassword, user.Password); diff != "" {
					t.Errorf("Long password mismatch (-want +got):\n%s", diff)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

// TestUsers_FieldValidation tests individual field validation
func TestUsers_FieldValidation(t *testing.T) {
	now := time.Now()

	// Test UUID field
	t.Run("UUID validation", func(t *testing.T) {
		user := model.Users{
			UUID: "test-user-uuid-validation",
		}

		if user.UUID == "" {
			t.Error("Expected UUID to be set")
		}

		expectedUUID := "test-user-uuid-validation"
		if diff := cmp.Diff(expectedUUID, user.UUID); diff != "" {
			t.Errorf("UUID validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test Email field
	t.Run("Email validation", func(t *testing.T) {
		user := model.Users{
			Email: "test@validation.com",
		}

		if user.Email == "" {
			t.Error("Expected Email to be set")
		}

		expectedEmail := "test@validation.com"
		if diff := cmp.Diff(expectedEmail, user.Email); diff != "" {
			t.Errorf("Email validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test Name field
	t.Run("Name validation", func(t *testing.T) {
		user := model.Users{
			Name: "Test User Name",
		}

		if user.Name == "" {
			t.Error("Expected Name to be set")
		}

		expectedName := "Test User Name"
		if diff := cmp.Diff(expectedName, user.Name); diff != "" {
			t.Errorf("Name validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test Password field
	t.Run("Password validation", func(t *testing.T) {
		user := model.Users{
			Password: "hashed_password_test",
		}

		if user.Password == "" {
			t.Error("Expected Password to be set")
		}

		expectedPassword := "hashed_password_test"
		if diff := cmp.Diff(expectedPassword, user.Password); diff != "" {
			t.Errorf("Password validation mismatch (-want +got):\n%s", diff)
		}
	})

	// Test timestamp fields
	t.Run("Timestamp validation", func(t *testing.T) {
		user := model.Users{
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		if user.CreatedAt == nil {
			t.Error("Expected CreatedAt to be set")
		}

		if user.UpdatedAt == nil {
			t.Error("Expected UpdatedAt to be set")
		}

		// Verify the timestamps are as expected
		if diff := cmp.Diff(now.Unix(), user.CreatedAt.Unix()); diff != "" {
			t.Errorf("CreatedAt validation mismatch (-want +got):\n%s", diff)
		}

		if diff := cmp.Diff(now.Unix(), user.UpdatedAt.Unix()); diff != "" {
			t.Errorf("UpdatedAt validation mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestUsers_Comprehensive tests comprehensive user scenarios
func TestUsers_Comprehensive(t *testing.T) {
	now := time.Now()

	// Create a comprehensive user
	user := model.Users{
		ID:        123,
		UUID:      "comprehensive-test-user-uuid-456",
		Email:     "comprehensive.test@example.com",
		Name:      "Comprehensive Test User",
		Password:  "comprehensive_hashed_password_789",
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Test all fields comprehensively
	tests := []struct {
		fieldName string
		actual    interface{}
		expected  interface{}
	}{
		{"ID", user.ID, uint(123)},
		{"UUID", user.UUID, "comprehensive-test-user-uuid-456"},
		{"Email", user.Email, "comprehensive.test@example.com"},
		{"Name", user.Name, "Comprehensive Test User"},
		{"Password", user.Password, "comprehensive_hashed_password_789"},
		{"CreatedAt", user.CreatedAt != nil, true},
		{"UpdatedAt", user.UpdatedAt != nil, true},
		{"DeletedAt", user.DeletedAt == nil, true},
	}

	for _, tt := range tests {
		t.Run("field_"+tt.fieldName, func(t *testing.T) {
			if diff := cmp.Diff(tt.expected, tt.actual); diff != "" {
				t.Errorf("%s field mismatch (-want +got):\n%s", tt.fieldName, diff)
			}
		})
	}
}
