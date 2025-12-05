package private

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type MemberRepository interface {
	GetMember(request request.MemberRequest) response.MemberResponse
	CreateMember(request request.MemberRequest) response.MemberResponse
	UpdateMember(request request.MemberRequest) response.MemberResponse
	DeleteMember(request request.MemberRequest) response.MemberResponse
}

type memberRepository struct {
	BaseConfig config.BaseConfig
}

func (rcvr memberRepository) GetMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/members"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) CreateMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/member"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) UpdateMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr memberRepository) DeleteMember(request request.MemberRequest) response.MemberResponse {
	var resp response.MemberResponse
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMemberRepository(conf config.BaseConfig) MemberRepository {
	return &memberRepository{BaseConfig: conf}
}
