package private

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Group interface {
	GetGroup(request request.Group) response.Groups
	CreateGroup(request request.Group) response.Groups
	UpdateGroup(request request.Group) response.Groups
	DeleteGroup(request request.Group) response.Groups
}

type group struct {
	BaseConfig config.BaseConfig
}

func (rcvr *group) GetGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/groups"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) CreateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/group"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) UpdateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) DeleteGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroup(conf config.BaseConfig) Group {
	return &group{BaseConfig: conf}
}
