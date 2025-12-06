package internal

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Member interface {
	BootstrapMemberForDB(request request.MemberRequest) response.MemberResponse
	GetMember(request request.MemberRequest) response.MemberResponse
	CreateMember(request request.MemberRequest) response.MemberResponse
	UpdateMember(request request.MemberRequest) response.MemberResponse
	DeleteMember(request request.MemberRequest) response.MemberResponse
}

type member struct {
	BaseConfig config.BaseConfig
}

func (rcvr *member) BootstrapMemberForDB(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
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

func (rcvr *member) GetMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/members"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) CreateMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/member"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) UpdateMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) DeleteMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMember(conf config.BaseConfig) Member {
	return &member{BaseConfig: conf}
}
