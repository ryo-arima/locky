package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockCommonRepository struct {
	invalidatedTokens map[string]bool
	jwtSecret         string
}

func (m *mockCommonRepository) GenerateJWTToken(email, role, uuid string) (string, string, error) {
	jti := "test-jti-" + email
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"uuid":  uuid,
		"jti":   jti,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", "", err
	}
	return tokenString, "refresh-token-" + email, nil
}

func (m *mockCommonRepository) ValidateJWTToken(tokenString string) (*middleware.Claims, error) {
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

func (m *mockCommonRepository) InvalidateToken(ctx context.Context, jti string, expiration time.Duration) error {
	if m.invalidatedTokens == nil {
		m.invalidatedTokens = make(map[string]bool)
	}
	m.invalidatedTokens[jti] = true
	return nil
}

func (m *mockCommonRepository) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	if m.invalidatedTokens == nil {
		return false, nil
	}
	return m.invalidatedTokens[jti], nil
}

func (m *mockCommonRepository) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func setupTestRouter() (*gin.Engine, *mock.MockUserRepository, *mockCommonRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	userRepo := &mock.MockUserRepository{
		Users: []model.Users{},
	}
	commonRepo := &mockCommonRepository{
		jwtSecret:         "test-secret",
		invalidatedTokens: make(map[string]bool),
	}

	// Create controller
	ctrl := controller.NewCommonControllerForPublic(userRepo, commonRepo)

	// Setup routes
	router.POST("/v1/share/common/auth/tokens", ctrl.Login)
	router.GET("/v1/share/common/auth/tokens/validate", ctrl.ValidateToken)
	router.POST("/v1/share/common/auth/tokens/refresh", ctrl.RefreshToken)
	router.DELETE("/v1/share/common/auth/tokens", ctrl.Logout)

	return router, userRepo, commonRepo
}

func TestCommonControllerForPublic_Login_Success(t *testing.T) {
	router, userRepo, _ := setupTestRouter()

	// Create test user with hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	userRepo.Users = append(userRepo.Users, model.Users{
		ID:       1,
		UUID:     "test-uuid-123",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: string(hashedPassword),
		Role:     "user",
	})

	// Prepare login request
	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginReq)

	// Make request
	req := httptest.NewRequest("POST", "/v1/share/common/auth/tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response, "access_token")
	assert.Contains(t, response, "refresh_token")
	assert.Equal(t, "SUCCESS", response["code"])
}

func TestCommonControllerForPublic_Login_InvalidCredentials(t *testing.T) {
	router, userRepo, _ := setupTestRouter()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	userRepo.Users = append(userRepo.Users, model.Users{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	})

	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/v1/share/common/auth/tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCommonControllerForPublic_Login_UserNotFound(t *testing.T) {
	router, _, _ := setupTestRouter()

	loginReq := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/v1/share/common/auth/tokens", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCommonControllerForPublic_ValidateToken_Success(t *testing.T) {
	router, _, commonRepo := setupTestRouter()

	// Generate token
	token, _, err := commonRepo.GenerateJWTToken("test@example.com", "user", "test-uuid")
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/v1/share/common/auth/tokens/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	user := response["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "user", user["role"])
}

func TestCommonControllerForPublic_ValidateToken_MissingHeader(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest("GET", "/v1/share/common/auth/tokens/validate", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCommonControllerForPublic_ValidateToken_InvalidToken(t *testing.T) {
	router, _, _ := setupTestRouter()

	req := httptest.NewRequest("GET", "/v1/share/common/auth/tokens/validate", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCommonControllerForPublic_Logout_Success(t *testing.T) {
	router, _, commonRepo := setupTestRouter()

	// Generate token
	token, _, err := commonRepo.GenerateJWTToken("test@example.com", "user", "test-uuid")
	require.NoError(t, err)

	req := httptest.NewRequest("DELETE", "/v1/share/common/auth/tokens", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewCommonControllerForPublic(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mockCommonRepository{jwtSecret: "test"}

	ctrl := controller.NewCommonControllerForPublic(userRepo, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForInternal(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mockCommonRepository{jwtSecret: "test"}

	ctrl := controller.NewCommonControllerForInternal(userRepo, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForPrivate(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mockCommonRepository{jwtSecret: "test"}

	ctrl := controller.NewCommonControllerForPrivate(userRepo, commonRepo)

	assert.NotNil(t, ctrl)
}

// Test repository initialization
func TestCommonRepositoryInitialization(t *testing.T) {
	conf := &config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret-key",
				},
			},
		},
	}

	repo := repository.NewCommonRepository(conf)
	assert.NotNil(t, repo)
}
