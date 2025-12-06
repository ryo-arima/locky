package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewUserControllerForPublic(t *testing.T) {
	cfg := config.BaseConfig{}
	userRepo := repository.NewUser(cfg)
	userUsecase := usecase.NewUser(userRepo)
	commonRepo := repository.NewCommon(cfg, nil)
	commonUsecase := usecase.NewCommon(commonRepo)

	ctrl := controller.NewUserPublic(userUsecase, commonUsecase, cfg)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForInternal(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	userUsecase := usecase.NewUser(userRepo)

	ctrl := controller.NewUserInternal(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewUserControllerForPrivate(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	userUsecase := usecase.NewUser(userRepo)
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewUserPrivate(userUsecase, commonRepo)

	assert.NotNil(t, ctrl)
}

// Test usecase initialization
func TestUserUsecaseInitialization(t *testing.T) {
	userRepo := &mock.MockUserRepository{}

	uc := usecase.NewUser(userRepo)

	assert.NotNil(t, uc)
}
