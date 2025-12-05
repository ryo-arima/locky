package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller/internal"
	"github.com/ryo-arima/locky/pkg/server/controller/private"
	"github.com/ryo-arima/locky/pkg/server/controller/public"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewUserControllerForPublic(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUserUsecase(userRepo)
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	conf := config.BaseConfig{}

	ctrl := public.NewUserController(userUsecase, commonRepo, conf)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForInternal(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	userUsecase := usecase.NewUserUsecase(userRepo)

	ctrl := internal.NewUserController(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForPrivate(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUserUsecase(userRepo)
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := private.NewUserController(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

// Test usecase initialization
func TestUserUsecaseInitialization(t *testing.T) {
	userRepo := &mock.MockUserRepository{}

	uc := usecase.NewUserUsecase(userRepo)

	assert.NotNil(t, uc)
}
