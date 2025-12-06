package internal

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Group interface {
	BootstrapGroupForDB(request request.Group) response.Groups
	GetGroup(request request.Group) response.Groups
	CreateGroup(request request.Group) response.Groups
	UpdateGroup(request request.Group) response.Groups
	DeleteGroup(request request.Group) response.Groups
}

type group struct {
	BaseConfig config.BaseConfig
}

func (rcvr *group) BootstrapGroupForDB(request request.Group) response.Groups {
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

func (rcvr *group) GetGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/groups"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) CreateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/group"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) UpdateGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *group) DeleteGroup(request request.Group) response.Groups {
	var resp response.Groups
	endpoint := fmt.Sprintf("%s/v1/internal/group/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_GROUP_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewGroup(conf config.BaseConfig) Group {
	return &group{BaseConfig: conf}
}
