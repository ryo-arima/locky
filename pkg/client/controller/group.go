package controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ryo-arima/locky/pkg/client/usecase"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/spf13/cobra"
)

func InitBootstrapGroupCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	bootstrapGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Initialize the groups table in the database.",
		Long:  "This command drops the existing groups table and recreates it based on the current model.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.Bootstrap(request.GroupRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return bootstrapGroupCmd
}

func InitCreateGroupCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	createGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Create a new group (internal).",
		Long:  "Creates a new group with the provided name.",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			out := uc.CreateInternal(request.GroupRequest{Name: name}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createGroupCmd.Flags().StringP("name", "n", "", "Group name (required)")
	createGroupCmd.MarkFlagRequired("name")
	return createGroupCmd
}

func InitCreateGroupCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	createGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Create a new group (admin).",
		Long:  "Creates a new group with the provided name.",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			out := uc.CreatePrivate(request.GroupRequest{Name: name}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createGroupCmd.Flags().StringP("name", "n", "", "Group name (required)")
	createGroupCmd.MarkFlagRequired("name")
	return createGroupCmd
}

func InitGetGroupCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	getGroupCmd := &cobra.Command{
		Use:     "groups",
		Aliases: []string{"group"},
		Short:   "Get a list of groups (internal).",
		Long:    "Retrieves a list of all groups visible to an authenticated app user.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetInternal(request.GroupRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getGroupCmd
}

func InitGetGroupCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	getGroupCmd := &cobra.Command{
		Use:     "groups",
		Aliases: []string{"group"},
		Short:   "Get a list of groups (admin).",
		Long:    "Retrieves a list of all groups visible to an admin.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetPrivate(request.GroupRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getGroupCmd
}

func InitUpdateGroupCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	updateGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Update a group (internal).",
		Long:  "Updates a group's name. Requires group ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}

			out := uc.UpdateInternal(request.GroupRequest{
				ID:   uint(id),
				Name: name,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateGroupCmd.Flags().StringP("id", "i", "", "Group ID to update (required)")
	updateGroupCmd.Flags().StringP("name", "n", "", "New group name (required)")
	updateGroupCmd.MarkFlagRequired("id")
	updateGroupCmd.MarkFlagRequired("name")
	return updateGroupCmd
}

func InitUpdateGroupCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	updateGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Update a group (admin).",
		Long:  "Updates a group's name. Requires group ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}

			out := uc.UpdatePrivate(request.GroupRequest{
				ID:   uint(id),
				Name: name,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateGroupCmd.Flags().StringP("id", "i", "", "Group ID to update (required)")
	updateGroupCmd.Flags().StringP("name", "n", "", "New group name (required)")
	updateGroupCmd.MarkFlagRequired("id")
	updateGroupCmd.MarkFlagRequired("name")
	return updateGroupCmd
}

func InitDeleteGroupCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	deleteGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Delete a group (internal).",
		Long:  "Deletes a group by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeleteInternal(request.GroupRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteGroupCmd.Flags().StringP("id", "i", "", "Group ID to delete (required)")
	deleteGroupCmd.MarkFlagRequired("id")
	return deleteGroupCmd
}

func InitDeleteGroupCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewGroupUsecase(conf)
	deleteGroupCmd := &cobra.Command{
		Use:   "group",
		Short: "Delete a group (admin).",
		Long:  "Deletes a group by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeletePrivate(request.GroupRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteGroupCmd.Flags().StringP("id", "i", "", "Group ID to delete (required)")
	deleteGroupCmd.MarkFlagRequired("id")
	return deleteGroupCmd
}
