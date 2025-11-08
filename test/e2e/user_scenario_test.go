package e2e
package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/client/usecase"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserScenarioTestSuite struct {
	suite.Suite
	config       config.BaseConfig
	userUsecase  usecase.UserUsecase
	adminToken   string
	testUserUUID string
}

func (suite *UserScenarioTestSuite) SetupSuite() {
	// Load configuration
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "etc/app.dev.yaml"
	}

	conf, err := config.NewBaseConfig(configFile)
	if err != nil {
		suite.T().Fatalf("Failed to load config: %v", err)
	}
	suite.config = conf

	// Initialize repositories and usecases
	userRepo := repository.NewUserRepository(suite.config)
	commonRepo := repository.NewCommonRepository(suite.config)
	suite.userUsecase = usecase.NewUserUsecase(userRepo, commonRepo)

	// Wait for server to be ready
	time.Sleep(2 * time.Second)
}

func (suite *UserScenarioTestSuite) TearDownSuite() {
	// Cleanup if needed
}

// Test Scenario 1: Anonymous user registration
func (suite *UserScenarioTestSuite) TestScenario01_AnonymousUserRegistration() {
	suite.Run("Create user as anonymous", func() {
		testEmail := fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8])
		testPassword := "SecurePassword123!"
		testName := "Test User"

		req := request.UserRequest{
			Email:    testEmail,
			Password: testPassword,
			Name:     testName,
		}

		resp := suite.userUsecase.CreateUserForPublic(req)
		
		assert.Equal(suite.T(), "SUCCESS", resp.Code, "User creation should succeed")
		assert.NotEmpty(suite.T(), resp.Users, "Response should contain user data")
		assert.Equal(suite.T(), testEmail, resp.Users[0].Email, "Email should match")
		assert.Equal(suite.T(), testName, resp.Users[0].Name, "Name should match")
		assert.NotEmpty(suite.T(), resp.Users[0].UUID, "UUID should be generated")
		
		suite.testUserUUID = resp.Users[0].UUID
	})
}

// Test Scenario 2: User authentication flow
func (suite *UserScenarioTestSuite) TestScenario02_UserAuthenticationFlow() {
	suite.Run("User login and token validation", func() {
		// First create a user
		testEmail := fmt.Sprintf("auth-test-%s@example.com", uuid.New().String()[:8])
		testPassword := "AuthTestPassword123!"
		
		createReq := request.UserRequest{
			Email:    testEmail,
			Password: testPassword,
			Name:     "Auth Test User",
		}
		
		createResp := suite.userUsecase.CreateUserForPublic(createReq)
		assert.Equal(suite.T(), "SUCCESS", createResp.Code, "User creation should succeed")

		// Login
		loginReq := request.CommonRequest{
			Email:    testEmail,
			Password: testPassword,
		}
		
		loginResp := suite.userUsecase.LoginUser(loginReq)
		assert.Equal(suite.T(), "SUCCESS", loginResp.Code, "Login should succeed")
		assert.NotEmpty(suite.T(), loginResp.Token, "Token should be returned")
		
		token := loginResp.Token

		// Validate token
		validateReq := request.CommonRequest{
			Token: token,
		}
		
		validateResp := suite.userUsecase.ValidateToken(validateReq)
		assert.Equal(suite.T(), "SUCCESS", validateResp.Code, "Token validation should succeed")

		// Get user info from token
		userInfoResp := suite.userUsecase.GetUserInfo(validateReq)
		assert.Equal(suite.T(), "SUCCESS", userInfoResp.Code, "Get user info should succeed")
		assert.Equal(suite.T(), testEmail, userInfoResp.Email, "Email should match")

		// Logout
		logoutResp := suite.userUsecase.LogoutUser(validateReq)
		assert.Equal(suite.T(), "SUCCESS", logoutResp.Code, "Logout should succeed")

		// Validate token after logout (should fail)
		invalidateResp := suite.userUsecase.ValidateToken(validateReq)
		assert.NotEqual(suite.T(), "SUCCESS", invalidateResp.Code, "Token validation should fail after logout")
	})
}

// Test Scenario 3: User CRUD operations with authentication
func (suite *UserScenarioTestSuite) TestScenario03_UserCRUDWithAuth() {
	suite.Run("Complete user CRUD flow with authentication", func() {
		// Create admin user and login
		adminEmail := fmt.Sprintf("admin-%s@example.com", uuid.New().String()[:8])
		adminPassword := "AdminPassword123!"
		
		createAdminReq := request.UserRequest{
			Email:    adminEmail,
			Password: adminPassword,
			Name:     "Admin User",
		}
		
		createAdminResp := suite.userUsecase.CreateUserForPublic(createAdminReq)
		assert.Equal(suite.T(), "SUCCESS", createAdminResp.Code, "Admin user creation should succeed")

		loginReq := request.CommonRequest{
			Email:    adminEmail,
			Password: adminPassword,
		}
		
		loginResp := suite.userUsecase.LoginUser(loginReq)
		assert.Equal(suite.T(), "SUCCESS", loginResp.Code, "Admin login should succeed")
		
		// Store token for subsequent requests
		suite.config.YamlConfig.Application.Client.Token = loginResp.Token

		// List users (internal endpoint)
		listReq := request.UserRequest{}
		listResp := suite.userUsecase.GetUserForInternal(listReq)
		assert.Equal(suite.T(), "SUCCESS", listResp.Code, "List users should succeed")
		assert.NotEmpty(suite.T(), listResp.Users, "User list should not be empty")

		// Create user via internal endpoint
		newUserEmail := fmt.Sprintf("internal-%s@example.com", uuid.New().String()[:8])
		newUserReq := request.UserRequest{
			Email:    newUserEmail,
			Password: "InternalPassword123!",
			Name:     "Internal Test User",
		}
		
		createResp := suite.userUsecase.CreateUserForInternal(newUserReq)
		assert.Equal(suite.T(), "SUCCESS", createResp.Code, "Create user via internal should succeed")
		
		createdUserUUID := createResp.Users[0].UUID

		// Update user
		updateReq := request.UserRequest{
			UUID:     createdUserUUID,
			Email:    newUserEmail,
			Password: "UpdatedPassword123!",
			Name:     "Updated Internal User",
		}
		
		updateResp := suite.userUsecase.UpdateUserForInternal(updateReq)
		assert.Equal(suite.T(), "SUCCESS", updateResp.Code, "Update user should succeed")
		assert.Equal(suite.T(), "Updated Internal User", updateResp.Users[0].Name, "Name should be updated")

		// Delete user
		deleteReq := request.UserRequest{
			UUID: createdUserUUID,
		}
		
		deleteResp := suite.userUsecase.DeleteUserForInternal(deleteReq)
		assert.Equal(suite.T(), "SUCCESS", deleteResp.Code, "Delete user should succeed")
	})
}

// Test Scenario 4: Private (admin) endpoints
func (suite *UserScenarioTestSuite) TestScenario04_PrivateAdminEndpoints() {
	suite.Run("Admin user management via private endpoints", func() {
		// Create and login as admin
		adminEmail := fmt.Sprintf("private-admin-%s@example.com", uuid.New().String()[:8])
		adminPassword := "PrivateAdminPassword123!"
		
		createAdminReq := request.UserRequest{
			Email:    adminEmail,
			Password: adminPassword,
			Name:     "Private Admin User",
		}
		
		suite.userUsecase.CreateUserForPublic(createAdminReq)

		loginReq := request.CommonRequest{
			Email:    adminEmail,
			Password: adminPassword,
		}
		
		loginResp := suite.userUsecase.LoginUser(loginReq)
		suite.config.YamlConfig.Application.Client.Token = loginResp.Token

		// Create user via private endpoint
		privateUserEmail := fmt.Sprintf("private-%s@example.com", uuid.New().String()[:8])
		privateUserReq := request.UserRequest{
			Email:    privateUserEmail,
			Password: "PrivatePassword123!",
			Name:     "Private Test User",
		}
		
		createResp := suite.userUsecase.CreateUserForPrivate(privateUserReq)
		assert.Equal(suite.T(), "SUCCESS", createResp.Code, "Create user via private should succeed")
		
		privateUserUUID := createResp.Users[0].UUID

		// List users via private endpoint
		listReq := request.UserRequest{}
		listResp := suite.userUsecase.GetUserForPrivate(listReq)
		assert.Equal(suite.T(), "SUCCESS", listResp.Code, "List users via private should succeed")

		// Update user via private endpoint
		updateReq := request.UserRequest{
			UUID:     privateUserUUID,
			Email:    privateUserEmail,
			Password: "UpdatedPrivatePassword123!",
			Name:     "Updated Private User",
		}
		
		updateResp := suite.userUsecase.UpdateUserForPrivate(updateReq)
		assert.Equal(suite.T(), "SUCCESS", updateResp.Code, "Update user via private should succeed")

		// Delete user via private endpoint
		deleteReq := request.UserRequest{
			UUID: privateUserUUID,
		}
		
		deleteResp := suite.userUsecase.DeleteUserForPrivate(deleteReq)
		assert.Equal(suite.T(), "SUCCESS", deleteResp.Code, "Delete user via private should succeed")
	})
}

// Test Scenario 5: Error handling and edge cases
func (suite *UserScenarioTestSuite) TestScenario05_ErrorHandlingAndEdgeCases() {
	suite.Run("Duplicate email registration", func() {
		testEmail := fmt.Sprintf("duplicate-%s@example.com", uuid.New().String()[:8])
		
		req := request.UserRequest{
			Email:    testEmail,
			Password: "DuplicatePassword123!",
			Name:     "Duplicate Test User",
		}
		
		// First creation should succeed
		resp1 := suite.userUsecase.CreateUserForPublic(req)
		assert.Equal(suite.T(), "SUCCESS", resp1.Code, "First user creation should succeed")

		// Second creation with same email should fail
		resp2 := suite.userUsecase.CreateUserForPublic(req)
		assert.NotEqual(suite.T(), "SUCCESS", resp2.Code, "Duplicate email should fail")
	})

	suite.Run("Invalid login credentials", func() {
		loginReq := request.CommonRequest{
			Email:    "nonexistent@example.com",
			Password: "WrongPassword123!",
		}
		
		loginResp := suite.userUsecase.LoginUser(loginReq)
		assert.NotEqual(suite.T(), "SUCCESS", loginResp.Code, "Login with invalid credentials should fail")
	})

	suite.Run("Invalid token validation", func() {
		validateReq := request.CommonRequest{
			Token: "invalid-token-string",
		}
		
		validateResp := suite.userUsecase.ValidateToken(validateReq)
		assert.NotEqual(suite.T(), "SUCCESS", validateResp.Code, "Validation with invalid token should fail")
	})

	suite.Run("Required fields validation", func() {
		// Missing email
		req1 := request.UserRequest{
			Password: "Password123!",
			Name:     "Test User",
		}
		
		resp1 := suite.userUsecase.CreateUserForPublic(req1)
		assert.NotEqual(suite.T(), "SUCCESS", resp1.Code, "Creation without email should fail")

		// Missing password
		req2 := request.UserRequest{
			Email: fmt.Sprintf("nopassword-%s@example.com", uuid.New().String()[:8]),
			Name:  "Test User",
		}
		
		resp2 := suite.userUsecase.CreateUserForPublic(req2)
		assert.NotEqual(suite.T(), "SUCCESS", resp2.Code, "Creation without password should fail")

		// Missing name
		req3 := request.UserRequest{
			Email:    fmt.Sprintf("noname-%s@example.com", uuid.New().String()[:8]),
			Password: "Password123!",
		}
		
		resp3 := suite.userUsecase.CreateUserForPublic(req3)
		assert.NotEqual(suite.T(), "SUCCESS", resp3.Code, "Creation without name should fail")
	})
}

func TestUserScenarioTestSuite(t *testing.T) {
	suite.Run(t, new(UserScenarioTestSuite))
}
