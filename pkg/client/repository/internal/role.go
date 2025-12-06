package internal

import (
	"net/http"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Role interface {
	ListRoles(filter share.RoleFilter) response.Roles
}

type role struct {
	base config.BaseConfig
}

func NewRole(base config.BaseConfig) Role {
	return &role{base: base}
}

func (r *role) endpoint(path string) string {
	return share.TrimEndpoint(r.base.YamlConfig.Application.Client.ServerEndpoint) + path
}

func (r *role) authReq(method, url string, body interface{}, out *response.Roles) error {
	return share.SendRequest(method, url, body, out)
}

func (r *role) ListRoles(filter share.RoleFilter) response.Roles {
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
