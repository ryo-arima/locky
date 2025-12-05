package private

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type GroupRepository interface {
	GetGroup(request request.GroupRequest) response.GroupResponse
	CreateGroup(request request.GroupRequest) response.GroupResponse
	UpdateGroup(request request.GroupRequest) response.GroupResponse
	DeleteGroup(request request.GroupRequest) response.GroupResponse
}

type groupRepository struct {
	BaseConfig config.BaseConfig
}

func (rcvr groupRepository) GetGroup(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/groups"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) CreateGroup(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/group"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) UpdateGroup(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) DeleteGroup(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroupRepository(conf config.BaseConfig) GroupRepository {
	return &groupRepository{BaseConfig: conf}
}
