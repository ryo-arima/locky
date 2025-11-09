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
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// TestMemberControllerForPublic_GetMembers tests the GetMembers endpoint
func TestMemberControllerForPublic_GetMembers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	memberRepo := repository.NewMemberRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	memberController := controller.NewMemberControllerForPublic(memberRepo, commonRepo)

	// Create request with empty JSON body
	req, _ := http.NewRequest("GET", "/api/public/members", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	memberController.GetMembers(c)

	// Expected response - GetMembers should succeed with empty JSON body
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

// TestMemberControllerForPublic_GetMembers_InvalidJSON tests GetMembers with invalid JSON
func TestMemberControllerForPublic_GetMembers_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	memberRepo := repository.NewMemberRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	memberController := controller.NewMemberControllerForPublic(memberRepo, commonRepo)

	// Create request with invalid JSON in body
	// Note: For GET requests, Gin's c.Bind() typically ignores request body and focuses on query parameters
	req, _ := http.NewRequest("GET", "/api/public/members", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	memberController.GetMembers(c)

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

// TestNewMemberControllerForPublic tests controller creation
func TestNewMemberControllerForPublic(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	memberRepo := repository.NewMemberRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	memberController := controller.NewMemberControllerForPublic(memberRepo, commonRepo)

	// Assert using go-cmp
	isNil := memberController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if memberController == nil {
		t.Error("Expected MemberControllerForPublic to be created, got nil")
	}
}

// TestMemberControllerForPublic_GetMembers_EmptyRequest tests with empty request
func TestMemberControllerForPublic_GetMembers_EmptyRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	memberRepo := repository.NewMemberRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	memberController := controller.NewMemberControllerForPublic(memberRepo, commonRepo)

	// Test data with empty request
	memberRequest := request.MemberRequest{}

	// Create request
	requestBody, _ := json.Marshal(memberRequest)
	req, _ := http.NewRequest("GET", "/api/public/members", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	memberController.GetMembers(c)

	// Expected response - should succeed with empty request
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}
}

// TestMemberControllerForPublic_TableDriven demonstrates table-driven tests
func TestMemberControllerForPublic_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		method       string
	}{
		{
			name:         "valid empty request",
			requestBody:  request.MemberRequest{},
			expectedCode: http.StatusOK,
			method:       "GET",
		},
		{
			name:         "nil request body",
			requestBody:  nil,
			expectedCode: http.StatusOK,
			method:       "GET",
		},
		{
			name:         "empty string request body",
			requestBody:  "",
			expectedCode: http.StatusOK,
			method:       "GET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			memberRepo := repository.NewMemberRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
			memberController := controller.NewMemberControllerForPublic(memberRepo, commonRepo)

			// Create request
			var requestBody []byte
			if tt.requestBody != nil && tt.requestBody != "" {
				requestBody, _ = json.Marshal(tt.requestBody)
			} else if tt.requestBody == "" {
				requestBody = []byte("")
			}

			req, _ := http.NewRequest(tt.method, "/api/public/members", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Test
			memberController.GetMembers(c)

			// Assert
			if diff := cmp.Diff(tt.expectedCode, w.Code); diff != "" {
				t.Errorf("Status code mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
