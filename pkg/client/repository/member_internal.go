package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type MemberInternal interface {
	BootstrapMemberForDB(request request.Member) response.Members
	GetMember(request request.Member) response.Members
	CreateMember(request request.Member) response.Members
	UpdateMember(request request.Member) response.Members
	DeleteMember(request request.Member) response.Members
}

type memberInternal struct {
	BaseConfig config.BaseConfig
}

func (rcvr *memberInternal) BootstrapMemberForDB(request request.Member) response.Members {
	var resp response.Members
	fmt.Println("BootstrapMemberForDB")

	if rcvr.BaseConfig.DBConnection == nil {
		if err := rcvr.BaseConfig.ConnectDB(); err != nil {
			resp.Code = "CLIENT_MEMBER_BOOTSTRAP_000"
			resp.Message = "Failed to connect database"
			return resp
		}
	}

	if rcvr.BaseConfig.DBConnection.Migrator().HasTable(&model.Members{}) {
		if err := rcvr.BaseConfig.DBConnection.Migrator().DropTable(&model.Members{}); err != nil {
			resp.Code = "CLIENT_MEMBER_BOOTSTRAP_001"
			resp.Message = fmt.Sprintf("Failed to drop existing table: %v", err)
			return resp
		}
	}

	if err := rcvr.BaseConfig.DBConnection.AutoMigrate(&model.Members{}); err != nil {
		resp.Code = "CLIENT_MEMBER_BOOTSTRAP_002"
		resp.Message = fmt.Sprintf("Failed to create Members table: %v", err)
		return resp
	}

	resp.Code = "SUCCESS"
	resp.Message = "Bootstrap for Member completed successfully"
	return resp
}

func (rcvr *memberInternal) GetMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/members"
	err := SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberInternal) CreateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/member"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberInternal) UpdateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberInternal) DeleteMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMemberInternal(conf config.BaseConfig) MemberInternal {
	return &memberInternal{BaseConfig: conf}
}
