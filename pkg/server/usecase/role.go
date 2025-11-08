package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type RoleUsecase interface {
	ListRoles(c *gin.Context) ([]string, error)
	GetRolePermissions(c *gin.Context, role string) ([]repository.RolePermission, error)
	CreateRole(c *gin.Context, role string, perms []repository.RolePermission) error
	UpdateRole(c *gin.Context, role string, perms []repository.RolePermission) error
	DeleteRole(c *gin.Context, role string) error
}

type roleUsecase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUsecase(roleRepo repository.RoleRepository) RoleUsecase {
	return &roleUsecase{
		roleRepo: roleRepo,
	}
}

func (uc *roleUsecase) ListRoles(c *gin.Context) ([]string, error) {
	return uc.roleRepo.ListRoles(c)
}

func (uc *roleUsecase) GetRolePermissions(c *gin.Context, role string) ([]repository.RolePermission, error) {
	return uc.roleRepo.GetRolePermissions(c, role)
}

func (uc *roleUsecase) CreateRole(c *gin.Context, role string, perms []repository.RolePermission) error {
	return uc.roleRepo.CreateRole(c, role, perms)
}

func (uc *roleUsecase) UpdateRole(c *gin.Context, role string, perms []repository.RolePermission) error {
	return uc.roleRepo.UpdateRole(c, role, perms)
}

func (uc *roleUsecase) DeleteRole(c *gin.Context, role string) error {
	return uc.roleRepo.DeleteRole(c, role)
}
