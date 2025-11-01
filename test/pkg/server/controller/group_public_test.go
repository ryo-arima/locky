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

// TestGroupControllerForPublic_GetGroups tests the GetGroups endpoint
func TestGroupControllerForPublic_GetGroups(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	groupRepo := repository.NewGroupRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	groupController := controller.NewGroupControllerForPublic(groupRepo, commonRepo)

	// Create request with empty JSON body
	req, _ := http.NewRequest("GET", "/api/public/groups", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	groupController.GetGroups(c)

	// Expected response - GetGroups should succeed with empty JSON body
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

// TestGroupControllerForPublic_GetGroups_InvalidJSON tests GetGroups with invalid JSON
func TestGroupControllerForPublic_GetGroups_InvalidJSON(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	groupRepo := repository.NewGroupRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	groupController := controller.NewGroupControllerForPublic(groupRepo, commonRepo)

	// Create request with invalid JSON in body
	// Note: For GET requests, Gin's c.Bind() typically ignores request body and focuses on query parameters
	req, _ := http.NewRequest("GET", "/api/public/groups", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	groupController.GetGroups(c)

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

// TestNewGroupControllerForPublic tests controller creation
func TestNewGroupControllerForPublic(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	groupRepo := repository.NewGroupRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	groupController := controller.NewGroupControllerForPublic(groupRepo, commonRepo)

	// Assert using go-cmp
	isNil := groupController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if groupController == nil {
		t.Error("Expected GroupControllerForPublic to be created, got nil")
	}
}

// TestGroupControllerForPublic_GetGroups_EmptyRequest tests with empty request
func TestGroupControllerForPublic_GetGroups_EmptyRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	groupRepo := repository.NewGroupRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)
	groupController := controller.NewGroupControllerForPublic(groupRepo, commonRepo)

	// Test data with empty request
	groupRequest := request.GroupRequest{}

	// Create request
	requestBody, _ := json.Marshal(groupRequest)
	req, _ := http.NewRequest("GET", "/api/public/groups", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Test
	groupController.GetGroups(c)

	// Expected response - should succeed with empty request
	expectedStatusCode := http.StatusOK

	// Assert using go-cmp for status code
	if diff := cmp.Diff(expectedStatusCode, w.Code); diff != "" {
		t.Errorf("Status code mismatch (-want +got):\n%s", diff)
	}
}
