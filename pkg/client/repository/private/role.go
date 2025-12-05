package private

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Role interface {
	ListRoles(filter repository.RoleFilter) response.RoleResponse
	CreateRole(req request.RolePermissionRequest) response.RoleResponse
	UpdateRole(role string, req request.RolePermissionRequest) response.RoleResponse
	DeleteRole(role string) response.RoleResponse
}

type role struct {
	base config.BaseConfig
}

func NewRole(base config.BaseConfig) Role {
	return &role{base: base}
}

func (r *role) endpoint(path string) string {
	return repository.TrimEndpoint(r.base.YamlConfig.Application.Client.ServerEndpoint) + path
}

func (r *role) authReq(method, url string, body interface{}, out *response.RoleResponse) error {
	return repository.SendRequest(method, url, body, out)
}

func (r *role) ListRoles(filter repository.RoleFilter) response.RoleResponse {
	url := r.endpoint("/v1/private/roles")
	if filter.ID != "" {
		url += "?id=" + filter.ID
	}
	var resp response.RoleResponse
	if err := r.authReq(http.MethodGet, url, nil, &resp); err != nil {
		resp.Code = "ROLE_LIST_ERROR"
		resp.Message = err.Error()
	}
	return resp
}

func (r *role) CreateRole(req request.RolePermissionRequest) response.RoleResponse {
	var resp response.RoleResponse
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

func (r *role) UpdateRole(role string, req request.RolePermissionRequest) response.RoleResponse {
	var resp response.RoleResponse
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

func (r *role) DeleteRole(role string) response.RoleResponse {
	var resp response.RoleResponse
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
