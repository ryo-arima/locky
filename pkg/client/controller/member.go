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

func InitBootstrapMemberCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	bootstrapMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Initialize the members table in the database.",
		Long:  "This command drops the existing members table and recreates it based on the current model.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.Bootstrap(request.MemberRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return bootstrapMemberCmd
}

func InitCreateMemberCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	createMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Create a new member association (internal).",
		Long:  "Creates a new member association between a user and a group.",
		Run: func(cmd *cobra.Command, args []string) {
			userUUID, _ := cmd.Flags().GetString("user-uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")
			out := uc.CreateInternal(request.MemberRequest{
				UserUUID:  userUUID,
				GroupUUID: groupUUID,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createMemberCmd.Flags().StringP("user-uuid", "u", "", "User UUID (required)")
	createMemberCmd.Flags().StringP("group-uuid", "g", "", "Group UUID (required)")
	createMemberCmd.MarkFlagRequired("user-uuid")
	createMemberCmd.MarkFlagRequired("group-uuid")
	return createMemberCmd
}

func InitCreateMemberCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	createMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Create a new member association (admin).",
		Long:  "Creates a new member association between a user and a group.",
		Run: func(cmd *cobra.Command, args []string) {
			userUUID, _ := cmd.Flags().GetString("user-uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")
			out := uc.CreatePrivate(request.MemberRequest{
				UserUUID:  userUUID,
				GroupUUID: groupUUID,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createMemberCmd.Flags().StringP("user-uuid", "u", "", "User UUID (required)")
	createMemberCmd.Flags().StringP("group-uuid", "g", "", "Group UUID (required)")
	createMemberCmd.MarkFlagRequired("user-uuid")
	createMemberCmd.MarkFlagRequired("group-uuid")
	return createMemberCmd
}

func InitGetMemberCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	getMemberCmd := &cobra.Command{
		Use:     "members",
		Aliases: []string{"member"},
		Short:   "Get a list of members (internal).",
		Long:    "Retrieves a list of all member associations visible to an authenticated app user.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetInternal(request.MemberRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getMemberCmd
}

func InitGetMemberCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	getMemberCmd := &cobra.Command{
		Use:     "members",
		Aliases: []string{"member"},
		Short:   "Get a list of members (admin).",
		Long:    "Retrieves a list of all member associations visible to an admin.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetPrivate(request.MemberRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getMemberCmd
}

func InitUpdateMemberCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	updateMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Update a member association (internal).",
		Long:  "Updates a member's user or group ID. Requires member ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			userUUID, _ := cmd.Flags().GetString("user-uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}

			out := uc.UpdateInternal(request.MemberRequest{
				ID:        uint(id),
				UserUUID:  userUUID,
				GroupUUID: groupUUID,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateMemberCmd.Flags().StringP("id", "i", "", "Member ID to update (required)")
	updateMemberCmd.Flags().StringP("user-uuid", "u", "", "New User UUID")
	updateMemberCmd.Flags().StringP("group-uuid", "g", "", "New Group UUID")
	updateMemberCmd.MarkFlagRequired("id")
	return updateMemberCmd
}

func InitUpdateMemberCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	updateMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Update a member association (admin).",
		Long:  "Updates a member's user or group ID. Requires member ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			userUUID, _ := cmd.Flags().GetString("user-uuid")
			groupUUID, _ := cmd.Flags().GetString("group-uuid")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}

			out := uc.UpdatePrivate(request.MemberRequest{
				ID:        uint(id),
				UserUUID:  userUUID,
				GroupUUID: groupUUID,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateMemberCmd.Flags().StringP("id", "i", "", "Member ID to update (required)")
	updateMemberCmd.Flags().StringP("user-uuid", "u", "", "New User UUID")
	updateMemberCmd.Flags().StringP("group-uuid", "g", "", "New Group UUID")
	updateMemberCmd.MarkFlagRequired("id")
	return updateMemberCmd
}

func InitDeleteMemberCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	deleteMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Delete a member association (internal).",
		Long:  "Deletes a member association by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeleteInternal(request.MemberRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteMemberCmd.Flags().StringP("id", "i", "", "Member ID to delete (required)")
	deleteMemberCmd.MarkFlagRequired("id")
	return deleteMemberCmd
}

func InitDeleteMemberCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewMemberUsecase(conf)
	deleteMemberCmd := &cobra.Command{
		Use:   "member",
		Short: "Delete a member association (admin).",
		Long:  "Deletes a member association by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeletePrivate(request.MemberRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteMemberCmd.Flags().StringP("id", "i", "", "Member ID to delete (required)")
	deleteMemberCmd.MarkFlagRequired("id")
	return deleteMemberCmd
}
