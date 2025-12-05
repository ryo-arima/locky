package controller

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

// RoleController: internal (read-only) role operations
// ListRoles also returns single details when specified with ?id=
// GetRole is deprecated as it's integrated into query-based approach
type RoleInternal interface {
	ListRoles(c *gin.Context)
}

type roleInternal struct {
	RoleUsecase usecase.Role
	enforcer    *casbin.Enforcer
}

func NewRoleInternal(roleUsecase usecase.Role, enf *casbin.Enforcer) RoleInternal {
	return &roleInternal{RoleUsecase: roleUsecase, enforcer: enf}
}

// ListRoles (internal) - read-only
func (rcvr *roleInternal) ListRoles(c *gin.Context) {
	if id := c.Query("id"); id != "" {
		perms, err := rcvr.RoleUsecase.GetRolePermissions(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_GET_ERROR", Message: err.Error(), Roles: []string{}})
			return
		}
		c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role permissions retrieved", Roles: []string{id}, Detail: perms})
		return
	}
	roles, err := rcvr.RoleUsecase.ListRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.RoleResponse{Code: "ROLE_LIST_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Roles retrieved", Roles: roles})
}
