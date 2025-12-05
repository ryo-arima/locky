package public

import (
	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserRepository interface {
	CreateUser(request request.UserRequest) response.UserResponse
}

type userRepository struct {
	BaseConfig config.BaseConfig
}

// CREATE
func (rcvr userRepository) CreateUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/public/user"
	err := repository.SendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PUBLIC_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUserRepository(conf config.BaseConfig) UserRepository {
	return &userRepository{BaseConfig: conf}
}
