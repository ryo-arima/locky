package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/server/controller/internal"
	"github.com/ryo-arima/locky/pkg/server/controller/private"
	"github.com/ryo-arima/locky/pkg/server/controller/public"
	mock "github.com/ryo-arima/locky/test/unit/mock/server"
	"github.com/stretchr/testify/assert"
)

func TestNewCommonControllerForPublic(t *testing.T) {
	userRepo := &mock.MockUserRepository{}
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := public.NewCommonController(userRepo, commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForInternal(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := internal.NewCommonController(commonRepo)

	assert.NotNil(t, ctrl)
}

func TestNewCommonControllerForPrivate(t *testing.T) {
	commonRepo := &mock.MockCommonRepository{JWTSecret: "test"}

	ctrl := private.NewCommonController(commonRepo)

	assert.NotNil(t, ctrl)
}
