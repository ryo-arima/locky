package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/server/controller/internal"
	"github.com/ryo-arima/locky/pkg/server/controller/private"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupControllerForInternal(t *testing.T) {
	groupRepo := &mock.MockGroupRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	ctrl := internal.NewGroupController(groupRepo, commonRepo)
	assert.NotNil(t, ctrl)
}

func TestNewGroupControllerForPrivate(t *testing.T) {
	groupRepo := &mock.MockGroupRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}
	ctrl := private.NewGroupController(groupRepo, commonRepo)
	assert.NotNil(t, ctrl)
}
