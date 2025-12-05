package internal

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type RoleRepository interface {
	ListRoles(filter repository.RoleFilter) response.RoleResponse
}

type roleRepository struct {
	base config.BaseConfig
}

func NewRoleRepository(base config.BaseConfig) RoleRepository {
	return &roleRepository{base: base}
}

func (r *roleRepository) endpoint(path string) string {
	return repository.TrimEndpoint(r.base.YamlConfig.Application.Client.ServerEndpoint) + path
}

func (r *roleRepository) authReq(method, url string, body interface{}, out *response.RoleResponse) error {
	return repository.SendRequest(method, url, body, out)
}

func (r *roleRepository) ListRoles(filter repository.RoleFilter) response.RoleResponse {
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
