package repository

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type RoleRepository interface {
	ListRolesInternal(filter RoleFilter) response.RoleResponse
	ListRolesPrivate(filter RoleFilter) response.RoleResponse
	CreateRole(req request.RolePermissionRequest) response.RoleResponse
	UpdateRole(role string, req request.RolePermissionRequest) response.RoleResponse
	DeleteRole(role string) response.RoleResponse
}

type roleRepository struct {
	base config.BaseConfig
}

type RoleFilter struct{ ID string }

func NewRoleRepository(base config.BaseConfig) RoleRepository { return &roleRepository{base: base} }

func (r *roleRepository) endpoint(path string) string {
	return strings.TrimRight(r.base.YamlConfig.Application.Client.ServerEndpoint, "/") + path
}

func (r *roleRepository) authReq(method, url string, body interface{}, out *response.RoleResponse) error {
	return sendRequest(method, url, body, out)
}

func (r *roleRepository) ListRolesInternal(filter RoleFilter) response.RoleResponse {
	url := r.endpoint("/v1/internal/roles")
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
func (r *roleRepository) ListRolesPrivate(filter RoleFilter) response.RoleResponse {
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
func (r *roleRepository) CreateRole(req request.RolePermissionRequest) response.RoleResponse {
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
func (r *roleRepository) UpdateRole(role string, req request.RolePermissionRequest) response.RoleResponse {
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
func (r *roleRepository) DeleteRole(role string) response.RoleResponse {
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

// Table formatting helper (used by usecase)
func rolesTableString(res response.RoleResponse) string {
	if res.Code != "SUCCESS" {
		return fmt.Sprintf("Code: %s\nMessage: %s\n", res.Code, res.Message)
	}
	// Single role retrieval (Detail contains permissions)
	if res.Detail != nil {
		// Format permissions display
		permLines := []string{}
		if perms, ok := res.Detail.([]interface{}); ok {
			for _, p := range perms {
				permLines = append(permLines, fmt.Sprintf("  - %v", p))
			}
		}
		return fmt.Sprintf("Role: %v\nPermissions:\n%v\n", res.Roles, strings.Join(permLines, "\n"))
	}
	// List: Sort and enumerate Roles alphabetically
	switch v := res.Roles.(type) {
	case []string:
		cp := make([]string, len(v))
		copy(cp, v)
		sort.Strings(cp)
		return fmt.Sprintf("Roles (%d): %s\n", len(cp), strings.Join(cp, ", "))
	default:
		return fmt.Sprintf("Roles: %v\n", v)
	}
}

// RolesTableStringAlias public function (for display use from other packages)
func RolesTableStringAlias(res response.RoleResponse) string { return rolesTableString(res) }
