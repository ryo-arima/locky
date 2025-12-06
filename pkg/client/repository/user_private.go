package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserPrivate interface {
	GetUser(request request.User) response.Users
	CreateUser(request request.User) response.Users
	UpdateUser(request request.User) response.Users
	DeleteUser(request request.User) response.Users
}

type userPrivate struct {
	BaseConfig config.BaseConfig
}

// GET
func (rcvr *userPrivate) GetUser(request request.User) response.Users {
	var resp response.Users
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/users"
	err := SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (rcvr *userPrivate) CreateUser(request request.User) response.Users {
	var resp response.Users
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/user"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr *userPrivate) UpdateUser(request request.User) response.Users {
	var resp response.Users
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr *userPrivate) DeleteUser(request request.User) response.Users {
	var resp response.Users
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUserPrivate(conf config.BaseConfig) UserPrivate {
	return &userPrivate{BaseConfig: conf}
}
