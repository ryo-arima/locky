package repository

import (
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserPublic interface {
	CreateUser(request request.User) response.Users
}

type userPublic struct {
	BaseConfig config.BaseConfig
}

// CREATE
func (rcvr *userPublic) CreateUser(request request.User) response.Users {
	var resp response.Users
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/public/user"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PUBLIC_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUserPublic(conf config.BaseConfig) UserPublic {
	return &userPublic{BaseConfig: conf}
}
