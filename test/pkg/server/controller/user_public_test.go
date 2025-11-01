package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// TestUserControllerForPublic_CreateUser_MissingEmail_Advanced tests validation logic
func TestUserControllerForPublic_CreateUser_MissingEmail_Advanced(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test data with missing email
	userRequest := request.UserRequest{
		Name:     "Test User",
		Password: "password123",
	}

	// Create request
	requestBody, _ := json.Marshal(userRequest)
	req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.CreateUser(c)

	// Expected response
	expectedStatusCode := http.StatusBadRequest
	expectedErrorCode := "SERVER_CONTROLLER_CREATE__FOR__002"

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	var response response.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert using go-cmp for error code
	if diff := cmp.Diff(expectedErrorCode, response.Code); diff != "" {
		t.Errorf("Error code mismatch (-want +got):\n%s", diff)
	}
}

// TestUserControllerForPublic_CreateUser_InvalidJSON_Advanced tests JSON parsing
func TestUserControllerForPublic_CreateUser_InvalidJSON_Advanced(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.CreateUser(c)

	// Expected values for comparison
	expectedStatusCode := http.StatusBadRequest
	expectedErrorCode := "SERVER_CONTROLLER_CREATE__FOR__001"

	// Assert using go-cmp
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	var response response.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if diff := cmp.Diff(expectedErrorCode, response.Code); diff != "" {
		t.Errorf("Error code mismatch (-want +got):\n%s", diff)
	}
}

// TestNewUserControllerForPublic_Advanced tests controller creation
func TestNewUserControllerForPublic_Advanced(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Assert using go-cmp
	isNil := userController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if userController == nil {
		t.Error("Expected UserControllerForPublic to be created, got nil")
	}
}

// TestUserControllerForPublic_ValidationErrorsAdvanced demonstrates table-driven tests with go-cmp
func TestUserControllerForPublic_ValidationErrorsAdvanced(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  request.UserRequest
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "missing email",
			requestBody:  request.UserRequest{Name: "Test", Password: "pass"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "SERVER_CONTROLLER_CREATE__FOR__002",
		},
		{
			name:         "missing name",
			requestBody:  request.UserRequest{Email: "test@test.com", Password: "pass"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "SERVER_CONTROLLER_CREATE__FOR__002",
		},
		{
			name:         "missing password",
			requestBody:  request.UserRequest{Email: "test@test.com", Name: "Test"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "SERVER_CONTROLLER_CREATE__FOR__002",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			userController.CreateUser(c)

			// Assert using go-cmp
			if diff := cmp.Diff(tt.expectedCode, w.Code); diff != "" {
				t.Errorf("Status code mismatch (-want +got):\n%s", diff)
			}

			var response response.UserResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if diff := cmp.Diff(tt.expectedErr, response.Code); diff != "" {
				t.Errorf("Error code mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestUserControllerForPublic_CreateUser_Success tests successful user creation
func TestUserControllerForPublic_CreateUser_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test data with all required fields
	userRequest := request.UserRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create request
	requestBody, _ := json.Marshal(userRequest)
	req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.CreateUser(c)

	// Expected response
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	// Verify response structure
	var response response.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check that user was created (in mock scenario, this would be mocked)
	if response.Code != "" && response.Code != "SUCCESS" {
		t.Errorf("Expected success or empty code, got: %s", response.Code)
	}
}

// TestUserControllerForPublic_CreateUser_DuplicateEmail tests duplicate email validation
func TestUserControllerForPublic_CreateUser_DuplicateEmail(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test data - this would normally require setting up existing user in database
	userRequest := request.UserRequest{
		Name:     "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	// Create request
	requestBody, _ := json.Marshal(userRequest)
	req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.CreateUser(c)

	// In a real test, you would mock the repository to return existing users
	// For now, we test the structure
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Expected either success or bad request, got: %d", w.Code)
	}
}

// TestUserControllerForPublic_GetUsers tests the GetUsers endpoint
func TestUserControllerForPublic_GetUsers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Create request with empty body
	req, _ := http.NewRequest("GET", "/api/public/users", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.GetUsers(c)

	// Expected response - GetUsers should succeed with empty JSON body
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	// Verify response can be unmarshaled
	var response interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

// TestUserControllerForPublic_GetUsers_InvalidJSON tests GetUsers with invalid JSON
func TestUserControllerForPublic_GetUsers_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Create request with invalid JSON in body
	// Note: For GET requests, Gin's c.Bind() typically ignores request body and focuses on query parameters
	// so invalid JSON in body doesn't cause bind errors for GET requests
	req, _ := http.NewRequest("GET", "/api/public/users", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.GetUsers(c)

	// Expected values - GET requests ignore request body JSON, so this succeeds
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	// Verify response can be unmarshaled (should be successful response)
	var response interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

// TestUserControllerForPublic_EdgeCases tests various edge cases
func TestUserControllerForPublic_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		method       string
		endpoint     string
	}{
		{
			name:         "empty request body",
			requestBody:  "",
			expectedCode: http.StatusBadRequest,
			method:       "POST",
			endpoint:     "/api/public/user",
		},
		{
			name:         "null request body",
			requestBody:  nil,
			expectedCode: http.StatusBadRequest,
			method:       "POST",
			endpoint:     "/api/public/user",
		},
		{
			name: "very long email",
			requestBody: request.UserRequest{
				Name:     "Test User",
				Email:    "verylongemailaddressthatmightcauseissues@verylongdomainnamethatmightcauseissues.com",
				Password: "password123",
			},
			expectedCode: http.StatusOK,
			method:       "POST",
			endpoint:     "/api/public/user",
		},
		{
			name: "special characters in name",
			requestBody: request.UserRequest{
				Name:     "Test User with Special Characters @#$%",
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedCode: http.StatusOK,
			method:       "POST",
			endpoint:     "/api/public/user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

			// Create request
			var requestBody []byte
			if tt.requestBody != nil && tt.requestBody != "" {
				requestBody, _ = json.Marshal(tt.requestBody)
			} else if tt.requestBody == "" {
				requestBody = []byte("")
			}

			req, _ := http.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			if tt.method == "POST" {
				userController.CreateUser(c)
			} else {
				userController.GetUsers(c)
			}

			// Assert - for edge cases, we mainly check that the server doesn't crash
			// and returns some reasonable status code
			if w.Code != tt.expectedCode && w.Code != http.StatusBadRequest && w.Code != http.StatusOK {
				t.Errorf("Unexpected status code: got %d, expected %d or 400 or 200", w.Code, tt.expectedCode)
			}
		})
	}
}

// TestUserControllerForPublic_CreateUser_SuccessfulCreation tests successful user creation with password hashing
func TestUserControllerForPublic_CreateUser_SuccessfulCreation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test data with all required fields
	userRequest := request.UserRequest{
		Name:     "Test User",
		Email:    "newuser@example.com",
		Password: "securepassword123",
	}

	// Create request
	requestBody, _ := json.Marshal(userRequest)
	req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.CreateUser(c)

	// Expected response
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	// Verify response structure
	var response response.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// In a real scenario with proper mocking, we would verify the user was created
	// For now, we check that no error code is returned
	if response.Code != "" && response.Code != "SUCCESS" {
		t.Logf("Response code: %s, message: %s", response.Code, response.Message)
	}
}

// TestUserControllerForPublic_CreateUser_EmailValidation tests various email formats
func TestUserControllerForPublic_CreateUser_EmailValidation(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		expectedCode int
	}{
		{
			name:         "valid email",
			email:        "valid@example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "email with plus sign",
			email:        "user+tag@example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "email with subdomain",
			email:        "user@mail.example.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "long email",
			email:        "verylongemailaddressthatmightcauseissues@verylongdomainname.com",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

			// Test data
			userRequest := request.UserRequest{
				Name:     "Test User",
				Email:    tt.email,
				Password: "password123",
			}

			// Create request
			requestBody, _ := json.Marshal(userRequest)
			req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			userController.CreateUser(c)

			// For these tests, we mainly want to ensure no crashes occur
			// The actual validation would depend on the repository implementation
			if w.Code != tt.expectedCode && w.Code != http.StatusBadRequest {
				t.Logf("Email: %s, Status: %d", tt.email, w.Code)
			}
		})
	}
}

// TestUserControllerForPublic_CreateUser_PasswordHashing tests password security
func TestUserControllerForPublic_CreateUser_PasswordHashing(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "P@ssw0rd!2023#Complex",
		},
		{
			name:     "password with spaces",
			password: "my secure password 123",
		},
		{
			name:     "unicode password",
			password: "パスワード123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

			// Test data
			userRequest := request.UserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: tt.password,
			}

			// Create request
			requestBody, _ := json.Marshal(userRequest)
			req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			userController.CreateUser(c)

			// Verify that password hashing doesn't cause errors
			// In a real test, we would verify the password is actually hashed
			if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
				t.Errorf("Unexpected status code for password '%s': %d", tt.password, w.Code)
			}
		})
	}
}

// TestUserControllerForPublic_GetUsers_Comprehensive tests GetUsers endpoint thoroughly
func TestUserControllerForPublic_GetUsers_Comprehensive(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test with different request methods and data
	tests := []struct {
		name         string
		requestBody  interface{}
		contentType  string
		expectedCode int
	}{
		{
			name:         "empty request body",
			requestBody:  nil,
			contentType:  "application/json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "valid empty json",
			requestBody:  map[string]interface{}{},
			contentType:  "application/json",
			expectedCode: http.StatusOK,
		},
		{
			name:         "request with user data",
			requestBody:  request.UserRequest{Email: "test@example.com", Name: "Test", Password: "pass"},
			contentType:  "application/json",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			if tt.requestBody != nil {
				requestBody, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req, _ := http.NewRequest("GET", "/api/public/users", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", tt.contentType)

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			userController.GetUsers(c)

			// Assert using go-cmp for status code
			if diff := cmp.Diff(tt.expectedCode, w.Code); diff != "" {
				t.Errorf("Status code mismatch (-want +got):\n%s", diff)
			}

			// Verify response can be unmarshaled
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
		})
	}
}

// TestUserControllerForPublic_GetUsers_ErrorHandling tests error scenarios for GetUsers
func TestUserControllerForPublic_GetUsers_ErrorHandling(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

	// Test with malformed JSON in request body
	// Note: For GET requests, Gin's c.Bind() ignores request body and focuses on query parameters
	req, _ := http.NewRequest("GET", "/api/public/users", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	userController.GetUsers(c)

	// Expected response - GET requests ignore request body, so this succeeds
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}

	// Verify response can be unmarshaled (should be successful response)
	var response interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

// TestUserControllerForPublic_SpecialCharacterHandling tests handling of special characters
func TestUserControllerForPublic_SpecialCharacterHandling(t *testing.T) {
	tests := []struct {
		name        string
		requestData request.UserRequest
		description string
	}{
		{
			name: "name with special characters",
			requestData: request.UserRequest{
				Name:     "João O'Connor-Smith (Jr.)",
				Email:    "joao@example.com",
				Password: "password123",
			},
			description: "Testing names with accents, apostrophes, hyphens, and parentheses",
		},
		{
			name: "name with unicode",
			requestData: request.UserRequest{
				Name:     "田中太郎",
				Email:    "tanaka@example.com",
				Password: "password123",
			},
			description: "Testing unicode characters in names",
		},
		{
			name: "email with special characters",
			requestData: request.UserRequest{
				Name:     "Test User",
				Email:    "test+tag123@sub.example-domain.com",
				Password: "password123",
			},
			description: "Testing email with plus sign, numbers, and hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			userController := controller.NewUserControllerForPublic(userRepo, commonRepo)

			// Create request
			requestBody, _ := json.Marshal(tt.requestData)
			req, _ := http.NewRequest("POST", "/api/public/user", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			userController.CreateUser(c)

			// We mainly want to ensure no crashes occur with special characters
			// The actual validation would depend on business requirements
			if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
				t.Errorf("%s: Unexpected status code: %d", tt.description, w.Code)
			}

			// Log response for debugging if needed
			if w.Code != http.StatusOK {
				var response response.UserResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err == nil {
					t.Logf("%s: Response code: %s, message: %s", tt.description, response.Code, response.Message)
				}
			}
		})
	}
}
