package controller

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// TestNewCommonControllerForPublic tests controller creation
func TestNewCommonControllerForPublic(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	userRepo := repository.NewUserRepository(testHelper.BaseConfig)
	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	commonController := controller.NewCommonControllerForPublic(userRepo, commonRepo)

	// Assert using go-cmp
	isNil := commonController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if commonController == nil {
		t.Error("Expected CommonControllerForPublic to be created, got nil")
	}
}

// TestNewCommonControllerForPrivate tests controller creation
func TestNewCommonControllerForPrivate(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	commonController := controller.NewCommonControllerForPrivate(commonRepo)

	// Assert using go-cmp
	isNil := commonController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if commonController == nil {
		t.Error("Expected CommonControllerForPrivate to be created, got nil")
	}
}

// TestNewCommonControllerForInternal tests controller creation
func TestNewCommonControllerForInternal(t *testing.T) {
	// Setup
	testHelper := NewTestHelper()
	defer testHelper.CleanupDB()

	commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

	// Test
	commonController := controller.NewCommonControllerForInternal(commonRepo)

	// Assert using go-cmp
	isNil := commonController == nil
	if diff := cmp.Diff(false, isNil); diff != "" {
		t.Errorf("Controller creation check mismatch (-want +got):\n%s", diff)
	}

	// Additional assertion for clarity
	if commonController == nil {
		t.Error("Expected CommonControllerForInternal to be created, got nil")
	}
}

// TestCommonControllers_TableDriven demonstrates table-driven tests for all common controllers
func TestCommonControllers_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		controllerType string
	}{
		{
			name:           "public common controller",
			controllerType: "public",
		},
		{
			name:           "private common controller",
			controllerType: "private",
		},
		{
			name:           "internal common controller",
			controllerType: "internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			testHelper := NewTestHelper()
			defer testHelper.CleanupDB()

			userRepo := repository.NewUserRepository(testHelper.BaseConfig)
			commonRepo := repository.NewCommonRepository(testHelper.BaseConfig, nil)

			// Test based on controller type
			var commonController interface{}
			switch tt.controllerType {
			case "public":
				commonController = controller.NewCommonControllerForPublic(userRepo, commonRepo)
			case "private":
				commonController = controller.NewCommonControllerForPrivate(commonRepo)
			case "internal":
				commonController = controller.NewCommonControllerForInternal(commonRepo)
			}

			// Assert
			if commonController == nil {
				t.Errorf("Expected %s controller to be created, got nil", tt.controllerType)
			}
		})
	}
}
