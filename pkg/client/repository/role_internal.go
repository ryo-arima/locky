package repository

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type RoleInternal interface {
	ListRoles(filter RoleFilter) response.Roles
}

type roleInternal struct {
	base config.BaseConfig
}

func NewRoleInternal(base config.BaseConfig) RoleInternal {
	return &roleInternal{base: base}
}

func (r *roleInternal) endpoint(path string) string {
	return TrimEndpoint(r.base.YamlConfig.Application.Client.ServerEndpoint) + path
}

func (r *roleInternal) authReq(method, url string, body interface{}, out *response.Roles) error {
	return SendRequest(method, url, body, out)
}

func (r *roleInternal) ListRoles(filter RoleFilter) response.Roles {
	url := r.endpoint("/v1/internal/roles")
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
