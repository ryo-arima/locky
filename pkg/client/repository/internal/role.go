package internal

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Role interface {
	ListRoles(filter repository.RoleFilter) response.RoleResponse
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
