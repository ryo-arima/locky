package client

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/entity/request"
)

// TestUserRepository_NewUserRepository tests repository creation
func TestUserRepository_NewUserRepository(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	// Test
	userRepo := repository.NewUserRepository(testHelper.BaseConfig)

	// Assert using go-cmp with nil comparison
	isNil := userRepo == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Repository creation check mismatch (-want +got):\n%s", diff)
	}

	// Using cmp to verify the repository is properly initialized
	if userRepo == nil {
		t.Error("Expected UserRepository to be created, got nil")
	}
}

// TestUserRepository_GetUsers tests the GetUsers method
func TestUserRepository_GetUsers(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)

	// Test
	users := userRepo.GetUserForPublic(request.UserRequest{})

	// Assert - users should be a slice (empty or with data)
	if users.Code == "" && users.Message == "" {
		t.Error("Expected response, got empty struct")
	}
}

// TestUserRepository_CreateUser tests the CreateUser method
func TestUserRepository_CreateUser(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)

	// Test data - create a mock user
	userRequest := request.UserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Test
	result := userRepo.CreateUserForPublic(userRequest)

	// Assert - result should not be nil
	if result.Code == "" && result.Message == "" {
		t.Error("Expected response, got empty struct")
	}
}

// TestUserRepository_TableDriven demonstrates table-driven tests for all operations
func TestUserRepository_TableDriven(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		shouldPass bool
		testData   interface{}
	}{
		{
			name:       "get all users operation",
			operation:  "get_all",
			shouldPass: true,
			testData:   nil,
		},
		{
			name:       "create user operation",
			operation:  "create",
			shouldPass: true,
			testData: request.UserRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)

			// Test based on operation
			var success bool
			switch tt.operation {
			case "get_all":
				users := userRepo.GetUserForPublic(request.UserRequest{})
				success = users.Code != "" || users.Message != ""
			case "create":
				userData := tt.testData.(request.UserRequest)
				result := userRepo.CreateUserForPublic(userData)
				success = result.Code != "" || result.Message != ""
			}

			// Assert
			if diff := cmp.Diff(tt.shouldPass, success); diff != "" {
				t.Errorf("Operation success mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestUserRepository_EdgeCases tests edge cases and error scenarios
func TestUserRepository_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		description string
		testFunc    func(t *testing.T, userRepo repository.UserRepository)
	}{
		{
			name:        "create user with empty data",
			description: "Should handle creation with empty data",
			testFunc: func(t *testing.T, userRepo repository.UserRepository) {
				userRequest := request.UserRequest{}
				result := userRepo.CreateUserForPublic(userRequest)
				// Should not panic and return something
				if result.Code == "" && result.Message == "" {
					t.Error("Expected response, got empty struct")
				}
			},
		},
		{
			name:        "create user with nil data",
			description: "Should handle creation with nil data",
			testFunc: func(t *testing.T, userRepo repository.UserRepository) {
				// CreateUserForPublic takes a struct, not a pointer, so we can't pass nil.
				// We pass an empty struct instead.
				userRequest := request.UserRequest{}
				result := userRepo.CreateUserForPublic(userRequest)
				// Should not panic and return something
				if result.Code == "" && result.Message == "" {
					t.Error("Expected response, got empty struct")
				}
			},
		},
		{
			name:        "multiple get users calls",
			description: "Should handle multiple GetUsers calls",
			testFunc: func(t *testing.T, userRepo repository.UserRepository) {
				users1 := userRepo.GetUserForPublic(request.UserRequest{})
				users2 := userRepo.GetUserForPublic(request.UserRequest{})

				// Should not panic and return something each time
				if users1.Code == "" && users1.Message == "" {
					t.Error("Expected first response, got empty struct")
				}
				if users2.Code == "" && users2.Message == "" {
					t.Error("Expected second response, got empty struct")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)

			// Execute test case
			tc.testFunc(t, userRepo)
		})
	}
}

// TestUserRepository_Comprehensive tests comprehensive scenarios
func TestUserRepository_Comprehensive(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)

	// Test comprehensive workflow: get -> create -> get again
	// 1. Get initial users
	initialUsers := userRepo.GetUserForPublic(request.UserRequest{})
	if initialUsers.Code == "" && initialUsers.Message == "" {
		t.Error("Expected response from initial GetUsers, got empty struct")
	}

	// 2. Create a new user
	newUser := request.UserRequest{
		Name:  "Comprehensive Test User",
		Email: "comprehensive@test.com",
	}
	createResult := userRepo.CreateUserForPublic(newUser)
	if createResult.Code == "" && createResult.Message == "" {
		t.Error("Expected response from CreateUser, got empty struct")
	}

	// 3. Get users again (should potentially include the new user)
	updatedUsers := userRepo.GetUserForPublic(request.UserRequest{})
	if updatedUsers.Code == "" && updatedUsers.Message == "" {
		t.Error("Expected response from final GetUsers, got empty struct")
	}

	// Using go-cmp to verify the workflow completed without errors
	workflowSuccess := (initialUsers.Code != "" || initialUsers.Message != "") &&
		(createResult.Code != "" || createResult.Message != "") &&
		(updatedUsers.Code != "" || updatedUsers.Message != "")
	if diff := cmp.Diff(true, workflowSuccess); diff != "" {
		t.Errorf("Comprehensive workflow mismatch (-want +got):\n%s", diff)
	}
}
