package repository

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type RolePrivate interface {
	ListRoles(filter RoleFilter) response.Roles
	CreateRole(req request.RolePermission) response.Roles
	UpdateRole(role string, req request.RolePermission) response.Roles
	DeleteRole(role string) response.Roles
}

type rolePrivate struct {
	base config.BaseConfig
}

func NewRolePrivate(base config.BaseConfig) RolePrivate {
	return &rolePrivate{base: base}
}

func (r *rolePrivate) endpoint(path string) string {
	return TrimEndpoint(r.base.YamlConfig.Application.Client.ServerEndpoint) + path
}

func (r *rolePrivate) authReq(method, url string, body interface{}, out *response.Roles) error {
	return SendRequest(method, url, body, out)
}

func (r *rolePrivate) ListRoles(filter RoleFilter) response.Roles {
	url := r.endpoint("/v1/private/roles")
	if filter.ID != "" {
		url += "?id=" + filter.ID
	}
	var resp response.Roles
	if err := r.authReq(http.MethodGet, url, nil, &resp); err != nil {
		resp.Code = "ROLE_LIST_ERROR"
		resp.Message = err.Error()
	}
	return resp
}

func (r *rolePrivate) CreateRole(req request.RolePermission) response.Roles {
	var resp response.Roles
	if req.Role == "" {
		resp.Code = "ROLE_CREATE_VALIDATION_ERROR"
		resp.Message = "role required"
		return resp
	}
	url := r.endpoint("/v1/private/role")
	if err := r.authReq(http.MethodPost, url, req, &resp); err != nil {
		resp.Code = "ROLE_CREATE_ERROR"
		resp.Message = err.Error()
	}
	return resp
}

func (r *rolePrivate) UpdateRole(role string, req request.RolePermission) response.Roles {
	var resp response.Roles
	if role == "" {
		resp.Code = "ROLE_UPDATE_VALIDATION_ERROR"
		resp.Message = "role id required"
		return resp
	}
	url := r.endpoint("/v1/private/role/" + role)
	if err := r.authReq(http.MethodPut, url, req, &resp); err != nil {
		resp.Code = "ROLE_UPDATE_ERROR"
		resp.Message = err.Error()
	}
	return resp
}

func (r *rolePrivate) DeleteRole(role string) response.Roles {
	var resp response.Roles
	if role == "" {
		resp.Code = "ROLE_DELETE_VALIDATION_ERROR"
		resp.Message = "role id required"
		return resp
	}
	url := r.endpoint("/v1/private/role/" + role)
	if err := r.authReq(http.MethodDelete, url, nil, &resp); err != nil {
		resp.Code = "ROLE_DELETE_ERROR"
		resp.Message = err.Error()
	}
	return resp
}
