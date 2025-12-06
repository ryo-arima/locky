package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/server/controller"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewCommonControllerForPublic(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewCommonPublic(userRepo, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForInternal(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewCommonInternal(commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForPrivate(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := controller.NewCommonPrivate(commonRepo)

	assert.NotNil(t, ctrl)
}
