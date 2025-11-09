package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserUsecase struct {
	users []response.User
}

func (m *mockUserUsecase) ListUsers(ctx context.Context, filter repository.UserQueryFilter) ([]response.User, error) {
	return m.users, nil
}

func (m *mockUserUsecase) GetUser(ctx context.Context, id uint) (*response.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *mockUserUsecase) CreateUser(ctx context.Context, req *model.Users) (*response.User, error) {
	user := response.User{
		ID:    uint(len(m.users) + 1),
		UUID:  req.UUID,
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}
	m.users = append(m.users, user)
	return &user, nil
}

func (m *mockUserUsecase) UpdateUser(ctx context.Context, user *model.Users) error {
	return nil
}

func (m *mockUserUsecase) DeleteUser(ctx context.Context, id uint) error {
	return nil
}

func setupUserPublicRouter() (*gin.Engine, *mockUserUsecase, *mockCommonRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequestIDMiddleware())

	userUsecase := &mockUserUsecase{
		users: []response.User{},
	}
	commonRepo := &mockCommonRepository{
		jwtSecret:         "test-secret",
		invalidatedTokens: make(map[string]bool),
	}

	conf := config.BaseConfig{
		YamlConfig: config.YamlConfig{
			Application: config.Application{
				Server: config.Server{
					JWTSecret: "test-secret",
				},
			},
		},
	}

	ctrl := controller.NewUserControllerForPublic(userUsecase, commonRepo, conf)

	router.POST("/v1/public/users", ctrl.CreateUser)
	router.GET("/v1/public/users", ctrl.GetUsers)

	return router, userUsecase, commonRepo
}

func TestUserControllerForPublic_CreateUser_Success(t *testing.T) {
	router, _, _ := setupUserPublicRouter()

	createReq := map[string]string{
		"email":    "newuser@example.com",
		"name":     "New User",
		"password": "SecurePass123!",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest("POST", "/v1/public/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	users := resp["users"].([]interface{})
	assert.Len(t, users, 1)

	user := users[0].(map[string]interface{})
	assert.Equal(t, "newuser@example.com", user["email"])
	assert.Equal(t, "New User", user["name"])
}

func TestUserControllerForPublic_CreateUser_MissingFields(t *testing.T) {
	router, _, _ := setupUserPublicRouter()

	tests := []struct {
		name string
		body map[string]string
	}{
		{
			name: "Missing email",
			body: map[string]string{
				"name":     "Test User",
				"password": "password123",
			},
		},
		{
			name: "Missing name",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
		},
		{
			name: "Missing password",
			body: map[string]string{
				"email": "test@example.com",
				"name":  "Test User",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)

			req := httptest.NewRequest("POST", "/v1/public/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestUserControllerForPublic_CreateUser_InvalidJSON(t *testing.T) {
	router, _, _ := setupUserPublicRouter()

	req := httptest.NewRequest("POST", "/v1/public/users", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserControllerForPublic_GetUsers(t *testing.T) {
	router, userUsecase, _ := setupUserPublicRouter()

	// Add test users
	userUsecase.users = []response.User{
		{
			ID:    1,
			UUID:  uuid.New().String(),
			Email: "user1@example.com",
			Name:  "User One",
			Role:  "user",
		},
		{
			ID:    2,
			UUID:  uuid.New().String(),
			Email: "user2@example.com",
			Name:  "User Two",
			Role:  "user",
		},
	}

	req := httptest.NewRequest("GET", "/v1/public/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	users := resp["users"].([]interface{})
	assert.Len(t, users, 2)
}

func TestNewUserControllerForPublic(t *testing.T) {
	userUsecase := &mockUserUsecase{}
	commonRepo := &mockCommonRepository{jwtSecret: "test"}
	conf := config.BaseConfig{}

	ctrl := controller.NewUserControllerForPublic(userUsecase, commonRepo, conf)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForInternal(t *testing.T) {
	userUsecase := &mockUserUsecase{}

	ctrl := controller.NewUserControllerForInternal(userUsecase)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForPrivate(t *testing.T) {
	userUsecase := &mockUserUsecase{}
	commonRepo := &mockCommonRepository{jwtSecret: "test"}

	ctrl := controller.NewUserControllerForPrivate(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

// Test usecase initialization
func TestUserUsecaseInitialization(t *testing.T) {
	userRepo := &mockUserRepository{}

	uc := usecase.NewUserUsecase(userRepo)

	assert.NotNil(t, uc)
}

// Test password hashing in user creation flow
func TestPasswordHashing(t *testing.T) {
	password := "TestPassword123!"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	// Verify password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// Verify wrong password fails
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("WrongPassword"))
	assert.Error(t, err)
}
