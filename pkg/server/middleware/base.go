package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

func ForPublic(conf config.BaseConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Public endpoints don't require authentication
		c.Next()
	}
}

func ForInternal(commonRepo repository.CommonRepository, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateJWTToken(c, commonRepo); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MIDDLEWARE_AUTH_001",
				"message": "Authentication required",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		// Previously: claims, _ := getUserFromContext(c) -> removed as unused
		c.Next()
	}
}

// ForPrivate: determine if email is included in admin.emails (not dependent solely on role claims)
func ForPrivate(commonRepo repository.CommonRepository, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := validateJWTToken(c, commonRepo); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MIDDLEWARE_AUTH_002",
				"message": "Admin authentication required",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		claims, _ := getUserFromContext(c)
		// Leave admin determination to Casbin policy
		_ = claims
		c.Next()
	}
}

// CasbinAuthorization: evaluate role(obj=resource, act=methodMapping) for each request
func CasbinAuthorization(enforcer *casbin.Enforcer, resource string, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := getUserFromContext(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"code": "MIDDLEWARE_AUTH_003", "message": "Authentication required"})
			c.Abort()
			return
		}
		allowed, err := enforcer.Enforce(claims.Role, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "MIDDLEWARE_AUTH_004", "message": "authorization error", "error": err.Error()})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"code": "MIDDLEWARE_AUTH_005", "message": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// validateJWTToken validates JWT token and sets user context
func validateJWTToken(c *gin.Context, commonRepo repository.CommonRepository) error {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("authorization header required")
	}

	// Check if it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("invalid authorization header format")
	}

	// Extract token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return errors.New("token required")
	}

	// First, parse the token to get the JTI without full validation
	unverifiedClaims, err := commonRepo.ParseTokenUnverified(tokenString)
	if err != nil {
		return fmt.Errorf("could not parse token: %w", err)
	}

	// Check if the token is in the denylist
	isInvalidated, err := commonRepo.IsTokenInvalidated(c.Request.Context(), unverifiedClaims.Jti)
	if err != nil {
		// If Redis check fails, deny access for security
		return fmt.Errorf("redis check failed: %w", err)
	}
	if isInvalidated {
		return errors.New("token has been invalidated")
	}

	// If not invalidated, proceed with full validation (signature, expiry)
	claims, err := commonRepo.ValidateJWTToken(tokenString)
	if err != nil {
		return err
	}

	// Set user context for use in controllers
	setUserContext(c, claims)
	return nil
}

// setUserContext sets user claims in gin context
func setUserContext(c *gin.Context, claims *model.JWTClaims) {
	c.Set("user_id", claims.UserID)
	c.Set("user_uuid", claims.UUID)
	c.Set("user_email", claims.Email)
	c.Set("user_name", claims.Name)
	c.Set("user_role", claims.Role)
	c.Set("user_claims", claims)
}

// getUserFromContext retrieves user claims from gin context
func getUserFromContext(c *gin.Context) (*model.JWTClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*model.JWTClaims)
	if !ok {
		return nil, false
	}

	return userClaims, true
}

// GetUserID is a helper function to get user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(uint)
	return id, ok
}

// GetUserUUID is a helper function to get user UUID from context
func GetUserUUID(c *gin.Context) (string, bool) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		return "", false
	}

	uuid, ok := userUUID.(string)
	return uuid, ok
}

// GetUserEmail is a helper function to get user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := userEmail.(string)
	return email, ok
}

// GetUserName is a helper function to get user name from context
func GetUserName(c *gin.Context) (string, bool) {
	userName, exists := c.Get("user_name")
	if !exists {
		return "", false
	}

	name, ok := userName.(string)
	return name, ok
}

// GetUserRole is a helper function to get user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := userRole.(string)
	return role, ok
}

// GetUserClaims is a helper function to get full user claims from context
func GetUserClaims(c *gin.Context) (*model.JWTClaims, bool) {
	return getUserFromContext(c)
}
