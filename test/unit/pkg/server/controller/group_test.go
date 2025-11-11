package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/server/controller"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupControllerForInternal(t *testing.T) {
	groupRepo := &mock.MockGroupRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	ctrl := controller.NewGroupControllerForInternal(groupRepo, commonRepo)
	assert.NotNil(t, ctrl)
}

func TestNewGroupControllerForPrivate(t *testing.T) {
	groupRepo := &mock.MockGroupRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	ctrl := controller.NewGroupControllerForPrivate(groupRepo, commonRepo)
	assert.NotNil(t, ctrl)
}
