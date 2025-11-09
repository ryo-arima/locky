package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCommonRepo struct {
	jwtSecret string
}

func (m *mockCommonRepo) ValidateJWTToken(tokenString string) (*middleware.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &middleware.Claims{
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
		UUID:  claims["uuid"].(string),
		Jti:   claims["jti"].(string),
	}, nil
}

func (m *mockCommonRepo) IsTokenInvalidated(tokenString string) (bool, error) {
	return false, nil
}

func generateTestToken(secret, email, role, uuid string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"uuid":  uuid,
		"jti":   "test-jti",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())

	router.GET("/test", func(c *gin.Context) {
		requestID := middleware.GetRequestID(c)
		assert.NotEmpty(t, requestID)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Test without middleware
	requestID := middleware.GetRequestID(c)
	assert.Empty(t, requestID)

	// Test with middleware
	c.Set("request_id", "test-request-id-123")
	requestID = middleware.GetRequestID(c)
	assert.Equal(t, "test-request-id-123", requestID)
}

func TestAuthenticationMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockCommonRepo{jwtSecret: "test-secret"}
	router.Use(middleware.AuthenticationMiddleware(mockRepo))

	router.GET("/protected", func(c *gin.Context) {
		email, _ := c.Get("email")
		role, _ := c.Get("role")
		uuid, _ := c.Get("uuid")

		c.JSON(http.StatusOK, gin.H{
			"email": email,
			"role":  role,
			"uuid":  uuid,
		})
	})

	// Generate valid token
	token, err := generateTestToken("test-secret", "test@example.com", "user", "test-uuid-123")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthenticationMiddleware_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockCommonRepo{jwtSecret: "test-secret"}
	router.Use(middleware.AuthenticationMiddleware(mockRepo))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticationMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mockCommonRepo{jwtSecret: "test-secret"}
	router.Use(middleware.AuthenticationMiddleware(mockRepo))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-string")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthorizationMiddleware_AdminAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	router.Use(middleware.AuthorizationMiddleware("admin"))

	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthorizationMiddleware_InsufficientPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	router.Use(middleware.AuthorizationMiddleware("admin"))

	router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestAuthorizationMiddleware_MultipleRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	router.Use(middleware.AuthorizationMiddleware("admin", "user"))

	router.GET("/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "access granted"})
	})

	req := httptest.NewRequest("GET", "/resource", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestCORSMiddleware_OptionsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.OPTIONS("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not reach here"})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Set user context
	c.Set("email", "test@example.com")
	c.Set("role", "admin")
	c.Set("uuid", "test-uuid-123")

	// Get user context
	email, role, uuid := middleware.GetUserContext(c)

	assert.Equal(t, "test@example.com", email)
	assert.Equal(t, "admin", role)
	assert.Equal(t, "test-uuid-123", uuid)
}

func TestGetUserContext_MissingValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// Get user context without setting values
	email, role, uuid := middleware.GetUserContext(c)

	assert.Empty(t, email)
	assert.Empty(t, role)
	assert.Empty(t, uuid)
}

func TestClaimsStructure(t *testing.T) {
	claims := middleware.Claims{
		Email: "test@example.com",
		Role:  "admin",
		UUID:  "test-uuid",
		Jti:   "test-jti",
	}

	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "admin", claims.Role)
	assert.Equal(t, "test-uuid", claims.UUID)
	assert.Equal(t, "test-jti", claims.Jti)
}
