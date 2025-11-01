package controller

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// RoleControllerForInternal: internal (read-only) role operations
// ListRoles also returns single details when specified with ?id=
// GetRole is deprecated as it's integrated into query-based approach
type RoleControllerForInternal interface {
	ListRoles(c *gin.Context)
}

type roleControllerForInternal struct {
	repo     repository.RoleRepository
	enforcer *casbin.Enforcer
}

func NewRoleControllerForInternal(repo repository.RoleRepository, enf *casbin.Enforcer) RoleControllerForInternal {
	return &roleControllerForInternal{repo: repo, enforcer: enf}
}

// ListRoles (internal) - read-only
func (rc *roleControllerForInternal) ListRoles(c *gin.Context) {
	if id := c.Query("id"); id != "" {
		perms, err := rc.repo.GetRolePermissions(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_GET_ERROR", Message: err.Error(), Roles: []string{}})
			return
		}
		c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role permissions retrieved", Roles: []string{id}, Detail: perms})
		return
	}
	roles, err := rc.repo.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.RoleResponse{Code: "ROLE_LIST_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Roles retrieved", Roles: roles})
}
