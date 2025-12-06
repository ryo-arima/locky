package private

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Member interface {
	GetMember(request request.Member) response.Members
	CreateMember(request request.Member) response.Members
	UpdateMember(request request.Member) response.Members
	DeleteMember(request request.Member) response.Members
}

type member struct {
	BaseConfig config.BaseConfig
}

func (rcvr *member) GetMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/members"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) CreateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/member"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) UpdateMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr *member) DeleteMember(request request.Member) response.Members {
	var resp response.Members
	endpoint := fmt.Sprintf("%s/v1/private/member/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_MEMBER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewMember(conf config.BaseConfig) Member {
	return &member{BaseConfig: conf}
}
