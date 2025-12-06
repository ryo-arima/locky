package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type MemberPrivate interface {
	GetMember(request request.Member) response.Members
	CreateMember(request request.Member) response.Members
	UpdateMember(request request.Member) response.Members
	DeleteMember(request request.Member) response.Members
}

type memberPrivate struct {
	BaseConfig config.BaseConfig
}

func (rcvr *memberPrivate) GetMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/members"
	err := SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberPrivate) CreateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/member"
	err := SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberPrivate) UpdateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *memberPrivate) DeleteMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMemberPrivate(conf config.BaseConfig) MemberPrivate {
	return &memberPrivate{BaseConfig: conf}
}
