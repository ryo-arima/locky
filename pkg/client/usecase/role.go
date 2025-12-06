package usecase

import (
	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
)

type Role interface {
	ListInternal(id, format string) string
	ListPrivate(id, format string) string
	Create(role string, perms []request.RolePermissionItem, format string) string
	Update(role string, perms []request.RolePermissionItem, format string) string
	Delete(role string, format string) string
}

type role struct{ repo repository.Role }

func NewRole(conf config.BaseConfig) Role {
	return &role{repo: repository.NewRole(conf)}
}

func (u *role) ListInternal(id, format string) string {
	resp := u.repo.ListRolesInternal(share.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *role) ListPrivate(id, format string) string {
	resp := u.repo.ListRolesPrivate(share.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *role) Create(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.repo.CreateRole(request.RolePermissionRequest{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *role) Update(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.repo.UpdateRole(role, request.RolePermissionRequest{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *role) Delete(role string, format string) string {
	resp := u.repo.DeleteRole(role)
	return Format(format, resp)
}
