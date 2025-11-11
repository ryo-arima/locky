package controller_test

import (
	"testing"

	"github.com/ryo-arima/locky/pkg/client/controller"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestInitBootstrapGroupCmdForAdminUser(t *testing.T) {
	conf := config.BaseConfig{}
	cmd := controller.InitBootstrapGroupCmdForAdminUser(conf)

	assert.NotNil(t, cmd)
	assert.Equal(t, "group", cmd.Use)
	assert.Contains(t, cmd.Short, "Initialize")
}

func TestInitCreateGroupCmdForAppUser(t *testing.T) {
	conf := config.BaseConfig{}
	cmd := controller.InitCreateGroupCmdForAppUser(conf)

	assert.NotNil(t, cmd)
	assert.Equal(t, "group", cmd.Use)
	assert.Contains(t, cmd.Short, "Create a new group")

	nameFlag := cmd.Flag("name")
	assert.NotNil(t, nameFlag)
}

func TestInitCreateGroupCmdForAdminUser(t *testing.T) {
	conf := config.BaseConfig{}
	cmd := controller.InitCreateGroupCmdForAdminUser(conf)

	assert.NotNil(t, cmd)
	assert.Equal(t, "group", cmd.Use)
	assert.Contains(t, cmd.Short, "Create a new group")

	nameFlag := cmd.Flag("name")
	assert.NotNil(t, nameFlag)
}
