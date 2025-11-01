package client

import (
	"github.com/ryo-arima/locky/pkg/client/controller"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/spf13/cobra"
)

type BaseCmdForAppUser struct {
	Create *cobra.Command
	Get    *cobra.Command
	Update *cobra.Command
	Delete *cobra.Command
}

func InitRootCmdForAppUser() *cobra.Command {
	var output string
	var rootCmdForAppUser = &cobra.Command{
		Use:   "locky-app",
		Short: "'locky' is a CLI tool to manage anniversaries",
		Long:  `''`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			controller.SetOutputFormat(output)
		},
	}
	rootCmdForAppUser.PersistentFlags().StringVarP(&output, "output", "o", "table", "Output format: table|json|yaml")
	return rootCmdForAppUser
}

func InitBaseCmdForAppUser() BaseCmdForAppUser {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create the value of a key",
		Long:  "create the value of a key",
	}
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get the value of a key",
		Long:  "get the value of a key",
	}
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update the value of a key",
		Long:  "update the value of a key",
	}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete the value of a key",
		Long:  "delete the value of a key",
	}
	baseCmdForAppUser := BaseCmdForAppUser{
		Create: createCmd,
		Get:    getCmd,
		Update: updateCmd,
		Delete: deleteCmd,
	}
	return baseCmdForAppUser
}

func ClientForAppUser(conf config.BaseConfig) {
	rootCmdForAppUser := InitRootCmdForAppUser()
	rootCmdForAppUser.CompletionOptions.HiddenDefaultCmd = true
	baseCmdForAppUser := InitBaseCmdForAppUser()

	// role (read-only) under get command
	baseCmdForAppUser.Get.AddCommand(controller.InitGetRoleCmdForApp(conf))

	//create
	createGroupCmdForAppUser := controller.InitCreateGroupCmdForAppUser(conf)
	baseCmdForAppUser.Create.AddCommand(createGroupCmdForAppUser)
	createMemberCmdForAppUser := controller.InitCreateMemberCmdForAppUser(conf)
	baseCmdForAppUser.Create.AddCommand(createMemberCmdForAppUser)
	rootCmdForAppUser.AddCommand(baseCmdForAppUser.Create)

	//get
	getUserCmdForAppUser := controller.InitGetUserCmdForAppUser(conf)
	baseCmdForAppUser.Get.AddCommand(getUserCmdForAppUser)
	getGroupCmdForAppUser := controller.InitGetGroupCmdForAppUser(conf)
	baseCmdForAppUser.Get.AddCommand(getGroupCmdForAppUser)
	getMemberCmdForAppUser := controller.InitGetMemberCmdForAppUser(conf)
	baseCmdForAppUser.Get.AddCommand(getMemberCmdForAppUser)
	rootCmdForAppUser.AddCommand(baseCmdForAppUser.Get)

	//update
	updateUserCmdForAppUser := controller.InitUpdateUserCmdForAppUser(conf)
	baseCmdForAppUser.Update.AddCommand(updateUserCmdForAppUser)
	updateGroupCmdForAppUser := controller.InitUpdateGroupCmdForAppUser(conf)
	baseCmdForAppUser.Update.AddCommand(updateGroupCmdForAppUser)
	updateMemberCmdForAppUser := controller.InitUpdateMemberCmdForAppUser(conf)
	baseCmdForAppUser.Update.AddCommand(updateMemberCmdForAppUser)
	rootCmdForAppUser.AddCommand(baseCmdForAppUser.Update)

	//delete
	deleteUserCmdForAppUser := controller.InitDeleteUserCmdForAppUser(conf)
	baseCmdForAppUser.Delete.AddCommand(deleteUserCmdForAppUser)
	deleteGroupCmdForAppUser := controller.InitDeleteGroupCmdForAppUser(conf)
	baseCmdForAppUser.Delete.AddCommand(deleteGroupCmdForAppUser)
	deleteMemberCmdForAppUser := controller.InitDeleteMemberCmdForAppUser(conf)
	baseCmdForAppUser.Delete.AddCommand(deleteMemberCmdForAppUser)
	rootCmdForAppUser.AddCommand(baseCmdForAppUser.Delete)

	rootCmdForAppUser.Execute()
}
