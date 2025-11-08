package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type MemberRepository interface {
	BootstrapMemberForDB(request request.MemberRequest) response.MemberResponse
	GetMemberForInternal(request request.MemberRequest) response.MemberResponse
	GetMemberForPrivate(request request.MemberRequest) response.MemberResponse
	CreateMemberForInternal(request request.MemberRequest) response.MemberResponse
	CreateMemberForPrivate(request request.MemberRequest) response.MemberResponse
	UpdateMemberForInternal(request request.MemberRequest) response.MemberResponse
	UpdateMemberForPrivate(request request.MemberRequest) response.MemberResponse
	DeleteMemberForInternal(request request.MemberRequest) response.MemberResponse
	DeleteMemberForPrivate(request request.MemberRequest) response.MemberResponse
}

type memberRepository struct {
	BaseConfig config.BaseConfig
}

// Bootstrap
func (rcvr memberRepository) BootstrapMemberForDB(request request.MemberRequest) response.MemberResponse {
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

// GET
func (rcvr memberRepository) GetMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/members"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) GetMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/members"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (rcvr memberRepository) CreateMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/member"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) CreateMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/member"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr memberRepository) UpdateMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) UpdateMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr memberRepository) DeleteMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) DeleteMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMemberRepository(conf config.BaseConfig) MemberRepository {
	return &memberRepository{BaseConfig: conf}
}
