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

func InitBootstrapUserCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	bootstrapUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Initialize the users table in the database.",
		Long:  "This command drops the existing users table and recreates it based on the current model.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.Bootstrap(request.UserRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return bootstrapUserCmd
}

func InitCreateUserCmdForAnonymousUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	createUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Create a new user (public registration).",
		Long:  "Creates a new user with the provided email, name, and password.",
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			name, _ := cmd.Flags().GetString("name")
			password, _ := cmd.Flags().GetString("password")

			out := uc.CreatePublic(request.UserRequest{
				Email:    email,
				Name:     name,
				Password: password,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createUserCmd.Flags().StringP("email", "e", "", "User email (required)")
	createUserCmd.Flags().StringP("name", "n", "", "User name (required)")
	createUserCmd.Flags().StringP("password", "p", "", "User password (required)")
	createUserCmd.MarkFlagRequired("email")
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("password")
	return createUserCmd
}

func InitCreateUserCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	createUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Create a new user (admin).",
		Long:  "Creates a new user with the provided email, name, and password.",
		Run: func(cmd *cobra.Command, args []string) {
			email, _ := cmd.Flags().GetString("email")
			name, _ := cmd.Flags().GetString("name")
			password, _ := cmd.Flags().GetString("password")

			out := uc.CreatePrivate(request.UserRequest{
				Email:    email,
				Name:     name,
				Password: password,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	createUserCmd.Flags().StringP("email", "e", "", "User email (required)")
	createUserCmd.Flags().StringP("name", "n", "", "User name (required)")
	createUserCmd.Flags().StringP("password", "p", "", "User password (required)")
	createUserCmd.MarkFlagRequired("email")
	createUserCmd.MarkFlagRequired("name")
	createUserCmd.MarkFlagRequired("password")
	return createUserCmd
}

func InitGetUserCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	getUserCmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Get a list of users (internal).",
		Long:    "Retrieves a list of all users visible to an authenticated app user.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetInternal(request.UserRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getUserCmd
}

func InitGetUserCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	getUserCmd := &cobra.Command{
		Use:     "users",
		Aliases: []string{"user"},
		Short:   "Get a list of users (admin).",
		Long:    "Retrieves a list of all users visible to an admin.",
		Run: func(cmd *cobra.Command, args []string) {
			out := uc.GetPrivate(request.UserRequest{}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	return getUserCmd
}

func InitUpdateUserCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	updateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Update a user (internal).",
		Long:  "Updates a user's details. Requires user ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			password, _ := cmd.Flags().GetString("password")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.UpdateInternal(request.UserRequest{
				ID:       uint(id),
				Name:     name,
				Password: password,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateUserCmd.Flags().StringP("id", "i", "", "User ID to update (required)")
	updateUserCmd.Flags().StringP("name", "n", "", "New user name")
	updateUserCmd.Flags().StringP("password", "p", "", "New user password")
	updateUserCmd.MarkFlagRequired("id")
	return updateUserCmd
}

func InitUpdateUserCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	updateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Update a user (admin).",
		Long:  "Updates a user's details. Requires user ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			password, _ := cmd.Flags().GetString("password")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.UpdatePrivate(request.UserRequest{
				ID:       uint(id),
				Name:     name,
				Password: password,
			}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	updateUserCmd.Flags().StringP("id", "i", "", "User ID to update (required)")
	updateUserCmd.Flags().StringP("name", "n", "", "New user name")
	updateUserCmd.Flags().StringP("password", "p", "", "New user password")
	updateUserCmd.MarkFlagRequired("id")
	return updateUserCmd
}

func InitDeleteUserCmdForAppUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	deleteUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Delete a user (internal).",
		Long:  "Deletes a user by their ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeleteInternal(request.UserRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteUserCmd.Flags().StringP("id", "i", "", "User ID to delete (required)")
	deleteUserCmd.MarkFlagRequired("id")
	return deleteUserCmd
}

func InitDeleteUserCmdForAdminUser(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewUserUsecase(conf)
	deleteUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Delete a user (admin).",
		Long:  "Deletes a user by their ID.",
		Run: func(cmd *cobra.Command, args []string) {
			idStr, _ := cmd.Flags().GetString("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Fatalf("Invalid ID: %v", err)
			}
			out := uc.DeletePrivate(request.UserRequest{ID: uint(id)}, GetOutputFormat())
			fmt.Print(out)
		},
	}
	deleteUserCmd.Flags().StringP("id", "i", "", "User ID to delete (required)")
	deleteUserCmd.MarkFlagRequired("id")
	return deleteUserCmd
}
