package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type GroupRepository interface {
	BootstrapGroupForDB(request request.GroupRequest) response.GroupResponse
	GetGroupForInternal(request request.GroupRequest) response.GroupResponse
	GetGroupForPrivate(request request.GroupRequest) response.GroupResponse
	CreateGroupForInternal(request request.GroupRequest) response.GroupResponse
	CreateGroupForPrivate(request request.GroupRequest) response.GroupResponse
	UpdateGroupForInternal(request request.GroupRequest) response.GroupResponse
	UpdateGroupForPrivate(request request.GroupRequest) response.GroupResponse
	DeleteGroupForInternal(request request.GroupRequest) response.GroupResponse
	DeleteGroupForPrivate(request request.GroupRequest) response.GroupResponse
}

type groupRepository struct {
	BaseConfig config.BaseConfig
}

// Bootstrap
func (rcvr groupRepository) BootstrapGroupForDB(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	fmt.Println("BootstrapGroupForDB")

	if rcvr.BaseConfig.DBConnection == nil {
		if err := rcvr.BaseConfig.ConnectDB(); err != nil {
			resp.Code = "CLIENT_GROUP_BOOTSTRAP_000"
			resp.Message = "Failed to connect database"
			return resp
		}
	}

	if rcvr.BaseConfig.DBConnection.Migrator().HasTable(&model.Groups{}) {
		if err := rcvr.BaseConfig.DBConnection.Migrator().DropTable(&model.Groups{}); err != nil {
			resp.Code = "CLIENT_GROUP_BOOTSTRAP_001"
			resp.Message = fmt.Sprintf("Failed to drop existing table: %v", err)
			return resp
		}
	}

	if err := rcvr.BaseConfig.DBConnection.AutoMigrate(&model.Groups{}); err != nil {
		resp.Code = "CLIENT_GROUP_BOOTSTRAP_002"
		resp.Message = fmt.Sprintf("Failed to create Groups table: %v", err)
		return resp
	}

	resp.Code = "SUCCESS"
	resp.Message = "Bootstrap for Group completed successfully"
	return resp
}

// GET
func (rcvr groupRepository) GetGroupForInternal(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/groups"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) GetGroupForPrivate(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/groups"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (rcvr groupRepository) CreateGroupForInternal(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/group"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) CreateGroupForPrivate(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/group"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr groupRepository) UpdateGroupForInternal(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) UpdateGroupForPrivate(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr groupRepository) DeleteGroupForInternal(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr groupRepository) DeleteGroupForPrivate(request request.GroupRequest) response.GroupResponse {
	var resp response.GroupResponse
	endpoint := fmt.Sprintf("%s/v1/private/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroupRepository(conf config.BaseConfig) GroupRepository {
	return &groupRepository{BaseConfig: conf}
}
