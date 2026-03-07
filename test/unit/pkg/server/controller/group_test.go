package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/controller"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewGroupControllerForInternal(t *testing.T) {
	cfg := config.BaseConfig{}
	groupRepo := repository.NewGroup(cfg)
	memberRepo := repository.NewMember(cfg)
	commonRepo := repository.NewCommon(cfg, nil)
	groupUsecase := usecase.NewGroup(groupRepo, memberRepo, nil)
	commonUsecase := usecase.NewCommon(commonRepo)
	ctrl := controller.NewGroupInternal(groupUsecase, commonUsecase)
	assert.NotNil(t, ctrl)
}

func TestNewGroupControllerForPrivate(t *testing.T) {
	cfg := config.BaseConfig{}
	groupRepo := repository.NewGroup(cfg)
	memberRepo := repository.NewMember(cfg)
	commonRepo := repository.NewCommon(cfg, nil)
	groupUsecase := usecase.NewGroup(groupRepo, memberRepo, nil)
	commonUsecase := usecase.NewCommon(commonRepo)
	ctrl := controller.NewGroupPrivate(groupUsecase, commonUsecase)
	assert.NotNil(t, ctrl)
}
