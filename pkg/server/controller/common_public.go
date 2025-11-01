// Package controller provides HTTP handlers for the Locky API.
//
// This package contains controllers for handling HTTP requests across different access levels:
//   - Public controllers: No authentication required
//   - Internal controllers: Authentication required, standard operations
//   - Private controllers: Admin authentication required
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// CommonControllerForPublic provides public authentication endpoints.
//
// This interface handles authentication operations that don't require
// prior authentication, following OpenStack Keystone design patterns.
//
// Available endpoints:
//   - ValidateToken: Validates JWT tokens and returns user information
//   - GetUserInfo: Returns user information from authenticated context
//   - Login: Handles user login and JWT token issuance
//   - RefreshToken: Refreshes JWT tokens using refresh tokens
//   - Logout: Handles user logout (token invalidation)
type CommonControllerForPublic interface {
	ValidateToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type commonControllerForPublic struct {
	UserRepository   repository.UserRepository
	CommonRepository repository.CommonRepository
}

// ValidateToken validates JWT token and returns user information.
//
// This endpoint validates a JWT token provided in the Authorization header
// and returns the user information contained within the token.
//
// Route: GET /v1/share/common/auth/tokens/validate
// Security: Bearer token required
//
// swagger:route GET /share/common/auth/tokens/validate Authentication validateToken
//
// # Validates JWT token and returns user information
//
// This endpoint validates a JWT token and returns the user information contained within it.
//
// Security:
//
//	bearerAuth: []
//
// Responses:
//
//	200: tokenValidationResponse
//	400: errorResponse
//	401: errorResponse
func (commonController commonControllerForPublic) ValidateToken(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "COMMON_VALIDATE_001",
			"message": "Authorization header required",
		})
		return
	}

	// Extract token
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Validate token
	claims, err := commonController.CommonRepository.ValidateJWTToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "COMMON_VALIDATE_002",
			"message": "Invalid token",
			"error":   err.Error(),
		})
		return
	}

	// Return user information from token
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Token is valid",
		"data": gin.H{
			"user_id":    claims.UserID,
			"uuid":       claims.UUID,
			"email":      claims.Email,
			"name":       claims.Name,
			"role":       claims.Role,
			"expires_at": claims.ExpiresAt,
		},
	})
}

// GetUserInfo returns current user information from JWT token.
//
// This endpoint retrieves user information from the JWT token context
// that was set by the authentication middleware.
//
// Route: GET /v1/share/common/auth/tokens/user
// Security: Bearer token required
//
// swagger:route GET /share/common/auth/tokens/user Authentication getUserInfo
//
// # Get current user information
//
// Returns user information extracted from the authenticated JWT token context.
//
// Security:
//
//	bearerAuth: []
//
// Responses:
//
//	200: userInfoResponse
//	401: errorResponse
func (commonController commonControllerForPublic) GetUserInfo(c *gin.Context) {
	// Get user claims from context (set by middleware)
	userClaims, exists := middleware.GetUserClaims(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    "COMMON_USER_INFO_001",
			"message": "User not authenticated",
		})
		return
	}

	// Return user information
	c.JSON(http.StatusOK, &response.CommonResponse{
		Code:    "SUCCESS",
		Message: "User information retrieved successfully",
		Commons: []response.Common{
			{
				ID:   userClaims.UserID,
				UUID: userClaims.UUID,
			},
		},
	})
}

// Login handles user login and returns JWT tokens.
//
// This endpoint authenticates users with email and password,
// then returns JWT access and refresh tokens upon successful authentication.
//
// Route: POST /v1/share/common/auth/tokens
// Security: No authentication required
//
// swagger:route POST /share/common/auth/tokens Authentication login
//
// # User login and token issuance
//
// Authenticates a user with email and password and returns JWT tokens.
//
// Responses:
//
//	200: loginResponse
//	400: errorResponse
//	401: errorResponse
//	500: errorResponse
func (commonController commonControllerForPublic) Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.LoginResponse{
			Code:    "AUTH_LOGIN_001",
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if loginRequest.Email == "" || loginRequest.Password == "" {
		c.JSON(http.StatusBadRequest, &response.LoginResponse{
			Code:    "AUTH_LOGIN_002",
			Message: "Email and password are required",
		})
		return
	}

	// Get all users to find matching email
	users := commonController.UserRepository.GetUsers()
	var foundUser *model.Users

	for _, user := range users {
		if user.Email == loginRequest.Email {
			foundUser = &user
			break
		}
	}

	if foundUser == nil {
		c.JSON(http.StatusUnauthorized, &response.LoginResponse{
			Code:    "AUTH_LOGIN_003",
			Message: "Invalid email or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, &response.LoginResponse{
			Code:    "AUTH_LOGIN_004",
			Message: "Invalid email or password",
		})
		return
	}

	// Determine user role (simple logic - can be enhanced)
	role := "user"
	baseConfig := commonController.CommonRepository.GetBaseConfig()
	for _, adminEmail := range baseConfig.YamlConfig.Application.Server.Admin.Emails {
		if foundUser.Email == adminEmail {
			role = "admin"
			break
		}
	}

	// Generate token pair
	tokenPair, err := commonController.CommonRepository.GenerateTokenPair(
		foundUser.ID,
		foundUser.UUID,
		foundUser.Email,
		foundUser.Name,
		role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.LoginResponse{
			Code:    "AUTH_LOGIN_005",
			Message: "Failed to generate tokens: " + err.Error(),
		})
		return
	}

	// Prepare user response
	userResponse := &response.User{
		ID:    foundUser.ID,
		UUID:  foundUser.UUID,
		Email: foundUser.Email,
		Name:  foundUser.Name,
	}

	c.JSON(http.StatusOK, &response.LoginResponse{
		Code:      "SUCCESS",
		Message:   "Login successful",
		TokenPair: tokenPair,
		User:      userResponse,
	})
}

// RefreshToken handles JWT token refresh.
//
// This endpoint accepts a valid refresh token and returns a new
// access token and refresh token pair.
//
// Route: POST /v1/share/common/auth/tokens/refresh
// Security: Refresh token required
//
// swagger:route POST /share/common/auth/tokens/refresh Authentication refreshToken
//
// # Refresh JWT token
//
// Exchanges a valid refresh token for new access and refresh tokens.
//
// Responses:
//
//	200: refreshTokenResponse
//	400: errorResponse
//	401: errorResponse
//	500: errorResponse
func (commonController commonControllerForPublic) RefreshToken(c *gin.Context) {
	var refreshRequest request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.RefreshTokenResponse{
			Code:    "AUTH_REFRESH_001",
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate refresh token
	claims, err := commonController.CommonRepository.ValidateJWTToken(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, &response.RefreshTokenResponse{
			Code:    "AUTH_REFRESH_002",
			Message: "Invalid refresh token: " + err.Error(),
		})
		return
	}

	// Generate new token pair
	tokenPair, err := commonController.CommonRepository.GenerateTokenPair(
		claims.UserID,
		claims.UUID,
		claims.Email,
		claims.Name,
		claims.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.RefreshTokenResponse{
			Code:    "AUTH_REFRESH_003",
			Message: "Failed to generate new tokens: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &response.RefreshTokenResponse{
		Code:      "SUCCESS",
		Message:   "Token refreshed successfully",
		TokenPair: tokenPair,
	})
}

// Logout handles user logout by invalidating the token.
//
// This endpoint invalidates a JWT token by adding it to a Redis denylist.
// The token will be rejected in subsequent validation attempts.
//
// Route: DELETE /v1/share/common/auth/tokens
// Security: Bearer token required
//
// swagger:route DELETE /share/common/auth/tokens Authentication logout
//
// # User logout and token invalidation
//
// Logs out the user by invalidating the JWT token and adding it to a denylist.
//
// Security:
//
//	bearerAuth: []
//
// Responses:
//
//	200: logoutResponse
//	400: errorResponse
//	500: errorResponse
func (commonController commonControllerForPublic) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "AUTH_LOGOUT_001",
			"message": "Authorization header required",
		})
		return
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "AUTH_LOGOUT_002",
			"message": "Token not found in Authorization header",
		})
		return
	}

	// Invalidate the token by adding it to the Redis denylist
	err := commonController.CommonRepository.InvalidateToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "AUTH_LOGOUT_003",
			"message": "Failed to invalidate token",
			"error":   err.Error(),
		})
		return
	}

	// Delete cache (interface method)
	commonController.CommonRepository.DeleteTokenCache(tokenString)

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Logout successful, token has been invalidated.",
	})
}

// NewCommonControllerForPublic creates a new instance of CommonControllerForPublic.
//
// This constructor function initializes a new CommonControllerForPublic with the
// required repository dependencies for handling authentication operations.
//
// Parameters:
//   - userRepository: Repository for user data operations
//   - commonRepository: Repository for common operations including JWT token management
//
// Returns:
//   - CommonControllerForPublic: Configured controller instance ready for use
func NewCommonControllerForPublic(userRepository repository.UserRepository, commonRepository repository.CommonRepository) CommonControllerForPublic {
	return &commonControllerForPublic{
		UserRepository:   userRepository,
		CommonRepository: commonRepository,
	}
}
