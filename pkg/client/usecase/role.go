package usecase

import (
	"github.com/ryo-arima/locky/pkg/client/repository"
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

type role struct {
	internalRepo repository.RoleInternal
	privateRepo  repository.RolePrivate
}

func NewRole(conf config.BaseConfig) Role {
	return &role{
		internalRepo: repository.NewRoleInternal(conf),
		privateRepo:  repository.NewRolePrivate(conf),
	}
}

func (u *role) ListInternal(id, format string) string {
	resp := u.internalRepo.ListRoles(repository.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *role) ListPrivate(id, format string) string {
	resp := u.privateRepo.ListRoles(repository.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *role) Create(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.privateRepo.CreateRole(request.RolePermission{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *role) Update(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.privateRepo.UpdateRole(role, request.RolePermission{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *role) Delete(role string, format string) string {
	resp := u.privateRepo.DeleteRole(role)
	return Format(format, resp)
}
