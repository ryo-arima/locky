package internal

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserRepository interface {
	BootstrapUserForDB(request request.UserRequest) response.UserResponse
	GetUser(request request.UserRequest) response.UserResponse
	UpdateUser(request request.UserRequest) response.UserResponse
	DeleteUser(request request.UserRequest) response.UserResponse
}

type userRepository struct {
	BaseConfig config.BaseConfig
}

// Bootstrap
func (rcvr userRepository) BootstrapUserForDB(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	fmt.Println("BootstrapUserForDB")

	if rcvr.BaseConfig.DBConnection == nil {
		if err := rcvr.BaseConfig.ConnectDB(); err != nil {
			resp.Code = "CLIENT_USER_BOOTSTRAP_000"
			resp.Message = "Failed to connect database"
			return resp
		}
	}

	if rcvr.BaseConfig.DBConnection.Migrator().HasTable(&model.Users{}) {
		if err := rcvr.BaseConfig.DBConnection.Migrator().DropTable(&model.Users{}); err != nil {
			resp.Code = "CLIENT_USER_BOOTSTRAP_001"
			resp.Message = fmt.Sprintf("Failed to drop existing table: %v", err)
			return resp
		}
	}

	if err := rcvr.BaseConfig.DBConnection.AutoMigrate(&model.Users{}); err != nil {
		resp.Code = "CLIENT_USER_BOOTSTRAP_002"
		resp.Message = fmt.Sprintf("Failed to create Users table: %v", err)
		return resp
	}

	resp.Code = "SUCCESS"
	resp.Message = "Bootstrap for User completed successfully"
	return resp
}

// GET
func (rcvr userRepository) GetUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/users"
	err := repository.SendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr userRepository) UpdateUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/internal/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr userRepository) DeleteUser(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/internal/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := repository.SendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUserRepository(conf config.BaseConfig) UserRepository {
	return &userRepository{BaseConfig: conf}
}
