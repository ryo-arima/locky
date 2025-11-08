package controller

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// RoleControllerForPrivate: administrative full CRUD (single retrieval via GET /roles?id=xxx)
type RoleControllerForPrivate interface {
	ListRoles(c *gin.Context)
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

type roleControllerForPrivate struct {
	repo     repository.RoleRepository
	enforcer *casbin.Enforcer
}

func NewRoleControllerForPrivate(repo repository.RoleRepository, enf *casbin.Enforcer) RoleControllerForPrivate {
	return &roleControllerForPrivate{repo: repo, enforcer: enf}
}

func (rc *roleControllerForPrivate) ListRoles(c *gin.Context) {
	if id := c.Query("id"); id != "" {
		perms, err := rc.repo.GetRolePermissions(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_GET_ERROR", Message: err.Error(), Roles: []string{}})
			return
		}
		c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role permissions retrieved", Roles: []string{id}, Detail: perms})
		return
	}
	roles, err := rc.repo.ListRoles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.RoleResponse{Code: "ROLE_LIST_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Roles retrieved", Roles: roles})
}

func (rc *roleControllerForPrivate) CreateRole(c *gin.Context) {
	var req request.RolePermissionRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_CREATE_BIND_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	if req.Role == "" {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_CREATE_VALIDATION_ERROR", Message: "role required", Roles: []string{}})
		return
	}
	perms := make([]repository.RolePermission, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		perms = append(perms, repository.RolePermission{Resource: p.Resource, Action: p.Action})
	}
	if err := rc.repo.CreateRole(c, req.Role, perms); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_CREATE_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role created", Roles: []string{req.Role}, Detail: perms})
}

func (rc *roleControllerForPrivate) UpdateRole(c *gin.Context) {
	role := c.Param("id")
	var req request.RolePermissionRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_UPDATE_BIND_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	if role == "" {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_UPDATE_VALIDATION_ERROR", Message: "role id(path) required", Roles: []string{}})
		return
	}
	perms := make([]repository.RolePermission, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		perms = append(perms, repository.RolePermission{Resource: p.Resource, Action: p.Action})
	}
	if err := rc.repo.UpdateRole(c, role, perms); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_UPDATE_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role updated", Roles: []string{role}, Detail: perms})
}

func (rc *roleControllerForPrivate) DeleteRole(c *gin.Context) {
	role := c.Param("id")
	if role == "" {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_DELETE_VALIDATION_ERROR", Message: "role id(path) required", Roles: []string{}})
		return
	}
	if err := rc.repo.DeleteRole(c, role); err != nil {
		c.JSON(http.StatusBadRequest, response.RoleResponse{Code: "ROLE_DELETE_ERROR", Message: err.Error(), Roles: []string{}})
		return
	}
	c.JSON(http.StatusOK, response.RoleResponse{Code: "SUCCESS", Message: "Role deleted", Roles: []string{role}})
}
