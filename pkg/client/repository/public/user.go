package public

import (
	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type User interface {
	CreateUser(request request.User) response.Users
}

type user struct {
	BaseConfig config.BaseConfig
}

// CREATE
func (rcvr *user) CreateUser(request request.User) response.Users {
	var resp response.Users
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/public/user"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PUBLIC_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUser(conf config.BaseConfig) User {
	return &user{BaseConfig: conf}
}
