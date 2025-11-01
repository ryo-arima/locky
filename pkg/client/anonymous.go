package client

import (
	"github.com/ryo-arima/locky/pkg/client/controller"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/spf13/cobra"
)

type BaseCmdForAnonymousUser struct {
	Get    *cobra.Command
	Create *cobra.Command
	Common *cobra.Command
}

func InitRootCmdForAnonymousUser() *cobra.Command {
	var output string
	var rootCmdForAnonymousUser = &cobra.Command{
		Use:   "locky-anonymous",
		Short: "'locky' is a CLI tool to manage anniversaries",
		Long:  `''`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			controller.SetOutputFormat(output)
		},
	}
	rootCmdForAnonymousUser.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table|json|yaml")
	return rootCmdForAnonymousUser
}

func InitBaseCmdForAnonymousUser() BaseCmdForAnonymousUser {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get the value of a key",
		Long:  "get the value of a key",
	}
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create the value of a key",
		Long:  "create the value of a key",
	}
	commonCmd := &cobra.Command{
		Use:   "common",
		Short: "common commands for anonymous users",
		Long:  "common commands for anonymous users",
	}
	baseCmdForAnonymousUser := BaseCmdForAnonymousUser{
		Get:    getCmd,
		Create: createCmd,
		Common: commonCmd,
	}
	return baseCmdForAnonymousUser
}

func ClientForAnonymousUser(conf config.BaseConfig) {
	rootCmdForAnonymousUser := InitRootCmdForAnonymousUser()
	rootCmdForAnonymousUser.CompletionOptions.HiddenDefaultCmd = true
	baseCmdForAnonymousUser := InitBaseCmdForAnonymousUser()

	//create
	createUserCmdForAnonymousUser := controller.InitCreateUserCmdForAnonymousUser(conf)
	baseCmdForAnonymousUser.Create.AddCommand(createUserCmdForAnonymousUser)
	rootCmdForAnonymousUser.AddCommand(baseCmdForAnonymousUser.Create)

	//common
	loginCmd := controller.InitCommonLoginCmd(conf)
	baseCmdForAnonymousUser.Common.AddCommand(loginCmd)
	validateCmd := controller.InitCommonValidateTokenCmd(conf)
	baseCmdForAnonymousUser.Common.AddCommand(validateCmd)
	userInfoCmd := controller.InitCommonUserInfoCmd(conf)
	baseCmdForAnonymousUser.Common.AddCommand(userInfoCmd)
	refreshCmd := controller.InitCommonRefreshTokenCmd(conf)
	baseCmdForAnonymousUser.Common.AddCommand(refreshCmd)
	logoutCmd := controller.InitCommonLogoutCmd(conf)
	baseCmdForAnonymousUser.Common.AddCommand(logoutCmd)
	rootCmdForAnonymousUser.AddCommand(baseCmdForAnonymousUser.Common)

	rootCmdForAnonymousUser.Execute()
}
