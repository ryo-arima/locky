package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type Role interface {
	ListRoles(c *gin.Context) ([]string, error)
	GetRolePermissions(c *gin.Context, role string) ([]repository.RolePermission, error)
	CreateRole(c *gin.Context, role string, perms []repository.RolePermission) error
	UpdateRole(c *gin.Context, role string, perms []repository.RolePermission) error
	DeleteRole(c *gin.Context, role string) error
}

type role struct {
	roleRepo repository.Role
}

func NewRole(roleRepo repository.Role) Role {
	return &role{
		roleRepo: roleRepo,
	}
}

func (uc *role) ListRoles(c *gin.Context) ([]string, error) {
	return uc.roleRepo.ListRoles(c)
}

func (uc *role) GetRolePermissions(c *gin.Context, role string) ([]repository.RolePermission, error) {
	return uc.roleRepo.GetRolePermissions(c, role)
}

func (uc *role) CreateRole(c *gin.Context, role string, perms []repository.RolePermission) error {
	return uc.roleRepo.CreateRole(c, role, perms)
}

func (uc *role) UpdateRole(c *gin.Context, role string, perms []repository.RolePermission) error {
	return uc.roleRepo.UpdateRole(c, role, perms)
}

func (uc *role) DeleteRole(c *gin.Context, role string) error {
	return uc.roleRepo.DeleteRole(c, role)
}
