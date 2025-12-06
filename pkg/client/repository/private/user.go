package private

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type User interface {
	GetUser(request request.UserRequest) response.UserResponse
	CreateUser(request request.UserRequest) response.UserResponse
	UpdateUser(request request.UserRequest) response.UserResponse
	DeleteUser(request request.UserRequest) response.UserResponse
}

type user struct {
	BaseConfig config.BaseConfig
}

// GET
func (rcvr *user) GetUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/users"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (rcvr *user) CreateUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/user"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr *user) UpdateUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr *user) DeleteUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUser(conf config.BaseConfig) User {
	return &user{BaseConfig: conf}
}
