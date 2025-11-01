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
func (mr memberRepository) BootstrapMemberForDB(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	fmt.Println("BootstrapMemberForDB")

	if mr.BaseConfig.DBConnection == nil {
		if err := mr.BaseConfig.ConnectDB(); err != nil {
			resp.Code = "CLIENT_MEMBER_BOOTSTRAP_000"
			resp.Message = "Failed to connect database"
			return resp
		}
	}

	if mr.BaseConfig.DBConnection.Migrator().HasTable(&model.Members{}) {
		if err := mr.BaseConfig.DBConnection.Migrator().DropTable(&model.Members{}); err != nil {
			resp.Code = "CLIENT_MEMBER_BOOTSTRAP_001"
			resp.Message = fmt.Sprintf("Failed to drop existing table: %v", err)
			return resp
		}
	}

	if err := mr.BaseConfig.DBConnection.AutoMigrate(&model.Members{}); err != nil {
		resp.Code = "CLIENT_MEMBER_BOOTSTRAP_002"
		resp.Message = fmt.Sprintf("Failed to create Members table: %v", err)
		return resp
	}

	resp.Code = "SUCCESS"
	resp.Message = "Bootstrap for Member completed successfully"
	return resp
}

// GET
func (mr memberRepository) GetMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/members"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (mr memberRepository) GetMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/members"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (mr memberRepository) CreateMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/member"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (mr memberRepository) CreateMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/member"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (mr memberRepository) UpdateMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (mr memberRepository) UpdateMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (mr memberRepository) DeleteMemberForInternal(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/internal/member/%d", mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (mr memberRepository) DeleteMemberForPrivate(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", mr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
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
