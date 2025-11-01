package usecase

import (
	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
)

type RoleUsecase interface {
	ListInternal(id, format string) string
	ListPrivate(id, format string) string
	Create(role string, perms []request.RolePermissionItem, format string) string
	Update(role string, perms []request.RolePermissionItem, format string) string
	Delete(role string, format string) string
}

type roleUsecase struct{ repo repository.RoleRepository }

func NewRoleUsecase(conf config.BaseConfig) RoleUsecase {
	return &roleUsecase{repo: repository.NewRoleRepository(conf)}
}

func (u *roleUsecase) ListInternal(id, format string) string {
	resp := u.repo.ListRolesInternal(repository.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *roleUsecase) ListPrivate(id, format string) string {
	resp := u.repo.ListRolesPrivate(repository.RoleFilter{ID: id})
	return Format(format, resp)
}
func (u *roleUsecase) Create(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.repo.CreateRole(request.RolePermissionRequest{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *roleUsecase) Update(role string, perms []request.RolePermissionItem, format string) string {
	resp := u.repo.UpdateRole(role, request.RolePermissionRequest{Role: role, Permissions: perms})
	return Format(format, resp)
}
func (u *roleUsecase) Delete(role string, format string) string {
	resp := u.repo.DeleteRole(role)
	return Format(format, resp)
}
