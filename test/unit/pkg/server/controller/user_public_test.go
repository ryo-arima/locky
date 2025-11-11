package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewUserControllerForPublic(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUserUsecase(userRepo)
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	conf := config.BaseConfig{}

	ctrl := controller.NewUserControllerForPublic(userUsecase, commonRepo, conf)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForInternal(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUserUsecase(userRepo)

	ctrl := controller.NewUserControllerForInternal(userUsecase)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForPrivate(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUserUsecase(userRepo)
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewUserControllerForPrivate(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

// Test usecase initialization
func TestUserUsecaseInitialization(t *testing.T) {
	userRepo := &mock.MockUserRepository{}

	uc := usecase.NewUserUsecase(userRepo)

	assert.NotNil(t, uc)
}
