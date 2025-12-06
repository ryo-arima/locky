package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type GroupPrivate interface {
	GetGroup(request request.Group) response.Groups
	CreateGroup(request request.Group) response.Groups
	UpdateGroup(request request.Group) response.Groups
	DeleteGroup(request request.Group) response.Groups
}

type groupPrivate struct {
	BaseConfig config.BaseConfig
}

func (rcvr *groupPrivate) GetGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/groups"
	err := SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupPrivate) CreateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/group"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupPrivate) UpdateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupPrivate) DeleteGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroupPrivate(conf config.BaseConfig) GroupPrivate {
	return &groupPrivate{BaseConfig: conf}
}
