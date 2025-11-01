package controller

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/usecase"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/spf13/cobra"
)

type permItems []request.RolePermissionItem

func (p *permItems) String() string { return fmt.Sprintf("%v", *p) }
func (p *permItems) Set(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		return fmt.Errorf("permission format resource:action")
	}
	*p = append(*p, request.RolePermissionItem{Resource: parts[0], Action: parts[1]})
	return nil
}
func (p *permItems) Type() string { return "perm" }

// Admin role get subcommand (under get to match other resources)
func InitGetRoleCmdForAdmin(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewRoleUsecase(conf)
	cmd := &cobra.Command{Use: "roles", Aliases: []string{"role"}, Short: "Get roles (admin)", Args: cobra.MaximumNArgs(1), Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) == 1 {
			id = args[0]
		}
		fmt.Print(uc.ListPrivate(id, GetOutputFormat()))
	}}
	return cmd
}

// Admin role create
func InitCreateRoleCmdForAdmin(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewRoleUsecase(conf)
	perms := permItems{}
	cmd := &cobra.Command{Use: "role", Short: "Create role (admin)", Args: cobra.ExactArgs(1), Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(uc.Create(args[0], perms, GetOutputFormat()))
	}}
	cmd.Flags().VarP(&perms, "perm", "p", "permission resource:action (repeatable)")
	return cmd
}

// Admin role update
func InitUpdateRoleCmdForAdmin(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewRoleUsecase(conf)
	perms := permItems{}
	cmd := &cobra.Command{Use: "role", Short: "Update role (admin)", Args: cobra.ExactArgs(1), Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(uc.Update(args[0], perms, GetOutputFormat()))
	}}
	cmd.Flags().VarP(&perms, "perm", "p", "permission resource:action (repeatable)")
	return cmd
}

// Admin role delete
func InitDeleteRoleCmdForAdmin(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewRoleUsecase(conf)
	cmd := &cobra.Command{Use: "role", Short: "Delete role (admin)", Args: cobra.ExactArgs(1), Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(uc.Delete(args[0], GetOutputFormat()))
	}}
	return cmd
}

// App (internal read-only)
func InitGetRoleCmdForApp(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewRoleUsecase(conf)
	cmd := &cobra.Command{Use: "roles", Aliases: []string{"role"}, Short: "Get roles (internal)", Args: cobra.MaximumNArgs(1), Run: func(cmd *cobra.Command, args []string) {
		id := ""
		if len(args) == 1 {
			id = args[0]
		}
		fmt.Print(uc.ListInternal(id, GetOutputFormat()))
	}}
	return cmd
}

// Compatibility wrappers for legacy InitRoleCmdForAdmin / InitRoleCmdForApp (can be removed in future)
func InitRoleCmdForAdmin(conf config.BaseConfig) *cobra.Command { return InitGetRoleCmdForAdmin(conf) }
func InitRoleCmdForApp(conf config.BaseConfig) *cobra.Command   { return InitGetRoleCmdForApp(conf) }
