package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type GroupInternal interface {
	BootstrapGroupForDB(request request.Group) response.Groups
	GetGroup(request request.Group) response.Groups
	CreateGroup(request request.Group) response.Groups
	UpdateGroup(request request.Group) response.Groups
	DeleteGroup(request request.Group) response.Groups
}

type groupInternal struct {
	BaseConfig config.BaseConfig
}

func (rcvr *groupInternal) BootstrapGroupForDB(request request.Group) response.Groups {
	var resp response.Groups
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

func (rcvr *groupInternal) GetGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/groups"
	err := SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupInternal) CreateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/group"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupInternal) UpdateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *groupInternal) DeleteGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroupInternal(conf config.BaseConfig) GroupInternal {
	return &groupInternal{BaseConfig: conf}
}
