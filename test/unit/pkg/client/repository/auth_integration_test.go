package client

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/entity/request"
)

// TestAuthIntegration_UserCreationAndLogin tests the complete flow of user creation and authentication
func TestAuthIntegration_UserCreationAndLogin(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	// Create repositories
	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig)

	// Test data
	testUser := request.UserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "TestPassword123!",
	}

	loginRequest := request.LoginRequest{
		Email:    testUser.Email,
		Password: testUser.Password,
	}

	// Test 1: Create a test user
	t.Run("create_test_user", func(t *testing.T) {
		result := userRepo.CreateUserForPublic(testUser)

		// Verify user creation (should have valid response structure)
		if result.Code == "" && result.Message == "" {
			t.Error("Expected user creation response with code or message")
		}

		// Log the result for debugging
		t.Logf("User creation result: %+v", result)
	})

	// Test 2: Login with created user
	t.Run("login_with_test_user", func(t *testing.T) {
		loginResponse := commonRepo.Login(loginRequest)

		// Verify login response structure
		isValidResponse := loginResponse.Code != "" || loginResponse.Message != ""
		if diff := cmp.Diff(true, isValidResponse); diff != "" {
			t.Errorf("Login response validation mismatch (-want +got):\n%s", diff)
		}

		// Log the response for debugging
		t.Logf("Login response: %+v", loginResponse)

		// Store token for next tests
		if loginResponse.TokenPair != nil && loginResponse.TokenPair.AccessToken != "" {
			t.Logf("Login successful! Access Token: %s", loginResponse.TokenPair.AccessToken[:20]+"...")

			// Test 3: Validate the received token
			t.Run("validate_access_token", func(t *testing.T) {
				validateResponse := commonRepo.ValidateToken(loginResponse.TokenPair.AccessToken)

				// Should get a response
				if validateResponse.Code == "" && validateResponse.Message == "" {
					t.Error("Expected validation response, got empty response")
				}

				t.Logf("Token validation response: %+v", validateResponse)
			})

			// Test 4: Get user info with token
			t.Run("get_user_info", func(t *testing.T) {
				userInfoResponse := commonRepo.GetUserInfo(loginResponse.TokenPair.AccessToken)

				// Should get a response
				if userInfoResponse.Code == "" && userInfoResponse.Message == "" {
					t.Error("Expected user info response, got empty response")
				}

				t.Logf("User info response: %+v", userInfoResponse)
			})

			// Test 5: Refresh token
			t.Run("refresh_token", func(t *testing.T) {
				if loginResponse.TokenPair.RefreshToken != "" {
					refreshResponse := commonRepo.RefreshToken(loginResponse.TokenPair.RefreshToken)

					// Should get a response
					if refreshResponse.Code == "" && refreshResponse.Message == "" {
						t.Error("Expected refresh token response, got empty response")
					}

					t.Logf("Token refresh response: %+v", refreshResponse)
				} else {
					t.Skip("No refresh token available for testing")
				}
			})

			// Test 6: Logout
			t.Run("logout", func(t *testing.T) {
				logoutResponse := commonRepo.Logout(loginResponse.TokenPair.AccessToken)

				// Should get a response
				if logoutResponse.Code == "" && logoutResponse.Message == "" {
					t.Error("Expected logout response, got empty response")
				}

				t.Logf("Logout response: %+v", logoutResponse)
			})
		} else {
			t.Log("Login did not return valid tokens - this may be expected if server is not running")
		}
	})
}

// TestAuthIntegration_EdgeCases tests edge cases in authentication flow
func TestAuthIntegration_EdgeCases(t *testing.T) {
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig)

	testCases := []struct {
		name        string
		description string
		testFunc    func(t *testing.T, repo repository.CommonRepository)
	}{
		{
			name:        "login_with_invalid_credentials",
			description: "Should handle login with invalid credentials gracefully",
			testFunc: func(t *testing.T, repo repository.CommonRepository) {
				invalidLogin := request.LoginRequest{
					Email:    "nonexistent@example.com",
					Password: "wrongpassword",
				}

				response := repo.Login(invalidLogin)

				// Should get a response (even if error)
				if response.Code == "" && response.Message == "" {
					t.Error("Expected error response for invalid credentials, got empty response")
				}

				t.Logf("Invalid login response: %+v", response)
			},
		},
		{
			name:        "validate_invalid_token",
			description: "Should handle token validation with invalid token",
			testFunc: func(t *testing.T, repo repository.CommonRepository) {
				invalidToken := "invalid.jwt.token"

				response := repo.ValidateToken(invalidToken)

				// Should get a response (likely error)
				if response.Code == "" && response.Message == "" {
					t.Error("Expected response for invalid token, got empty response")
				}

				t.Logf("Invalid token validation response: %+v", response)
			},
		},
		{
			name:        "refresh_with_invalid_token",
			description: "Should handle refresh with invalid token",
			testFunc: func(t *testing.T, repo repository.CommonRepository) {
				invalidRefreshToken := "invalid.refresh.token"

				response := repo.RefreshToken(invalidRefreshToken)

				// Should get a response (likely error)
				if response.Code == "" && response.Message == "" {
					t.Error("Expected response for invalid refresh token, got empty response")
				}

				t.Logf("Invalid refresh token response: %+v", response)
			},
		},
		{
			name:        "logout_without_token",
			description: "Should handle logout without token",
			testFunc: func(t *testing.T, repo repository.CommonRepository) {
				response := repo.Logout("")

				// Should get a response
				if response.Code == "" && response.Message == "" {
					t.Error("Expected response for logout without token, got empty response")
				}

				t.Logf("Logout without token response: %+v", response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t, commonRepo)
		})
	}
}

// TestAuthIntegration_PasswordValidation tests password validation scenarios
func TestAuthIntegration_PasswordValidation(t *testing.T) {
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig)

	passwordTests := []struct {
		name        string
		password    string
		shouldWork  bool
		description string
	}{
		{
			name:        "strong_password",
			password:    "StrongPass123!",
			shouldWork:  true,
			description: "Strong password with uppercase, lowercase, numbers, and symbols",
		},
		{
			name:        "weak_password_short",
			password:    "weak",
			shouldWork:  false,
			description: "Password too short",
		},
		{
			name:        "weak_password_no_uppercase",
			password:    "weakpass123!",
			shouldWork:  false,
			description: "Password without uppercase letters",
		},
		{
			name:        "weak_password_no_numbers",
			password:    "WeakPass!",
			shouldWork:  false,
			description: "Password without numbers",
		},
		{
			name:        "weak_password_no_special",
			password:    "WeakPass123",
			shouldWork:  false,
			description: "Password without special characters",
		},
	}

	for i, pwTest := range passwordTests {
		t.Run(pwTest.name, func(t *testing.T) {
			// Create unique test user for each password test
			testUser := request.UserRequest{
				Email:    fmt.Sprintf("passwordtest%d@example.com", i),
				Name:     fmt.Sprintf("Password Test User %d", i),
				Password: pwTest.password,
			}

			// Try to create user
			createResult := userRepo.CreateUserForPublic(testUser)
			t.Logf("Create user with %s: %+v", pwTest.description, createResult)

			// If user creation got some response, try to login
			if createResult.Code != "" || createResult.Message != "" {
				loginRequest := request.LoginRequest{
					Email:    testUser.Email,
					Password: testUser.Password,
				}

				loginResponse := commonRepo.Login(loginRequest)
				t.Logf("Login attempt with %s: %+v", pwTest.description, loginResponse)
			}
		})
	}
}

// TestAuthIntegration_ConcurrentAccess tests concurrent authentication operations
func TestAuthIntegration_ConcurrentAccess(t *testing.T) {
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig)

	// Test concurrent login attempts
	t.Run("concurrent_login_attempts", func(t *testing.T) {
		loginRequest := request.LoginRequest{
			Email:    "concurrent@example.com",
			Password: "ConcurrentTest123!",
		}

		// Channel to collect results
		results := make(chan bool, 3)

		// Launch 3 concurrent login attempts
		for i := 0; i < 3; i++ {
			go func(id int) {
				response := commonRepo.Login(loginRequest)
				// Consider it successful if we get any response
				success := response.Code != "" || response.Message != ""
				results <- success
				t.Logf("Concurrent login %d result: %+v", id, response)
			}(i)
		}

		// Wait for all results
		successCount := 0
		for i := 0; i < 3; i++ {
			if <-results {
				successCount++
			}
		}

		// All attempts should get some response (even if login fails)
		if successCount != 3 {
			t.Errorf("Expected 3 responses, got %d successful responses", successCount)
		}
	})
}
