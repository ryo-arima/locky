package repository

import (
	"fmt"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserRepository interface {
	BootstrapUserForDB(request request.UserRequest) response.UserResponse
	GetUserForInternal(request request.UserRequest) response.UserResponse
	GetUserForPrivate(request request.UserRequest) response.UserResponse
	CreateUserForPublic(request request.UserRequest) response.UserResponse
	CreateUserForPrivate(request request.UserRequest) response.UserResponse
	UpdateUserForInternal(request request.UserRequest) response.UserResponse
	UpdateUserForPrivate(request request.UserRequest) response.UserResponse
	DeleteUserForInternal(request request.UserRequest) response.UserResponse
	DeleteUserForPrivate(request request.UserRequest) response.UserResponse
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
func (rcvr userRepository) GetUserForInternal(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/internal/users"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_GET_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr userRepository) GetUserForPrivate(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/users"
	err := sendRequest("GET", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_GET_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// CREATE
func (rcvr userRepository) CreateUserForPublic(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/public/user"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PUBLIC_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr userRepository) CreateUserForPrivate(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/private/user"
	err := sendRequest("POST", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_CREATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// UPDATE
func (rcvr userRepository) UpdateUserForInternal(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/internal/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_UPDATE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr userRepository) UpdateUserForPrivate(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("PUT", endpoint, request, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_UPDATE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

// DELETE
func (rcvr userRepository) DeleteUserForInternal(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/internal/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_DELETE_INTERNAL_001"
		resp.Message = err.Error()
	}
	return resp
}

func (rcvr userRepository) DeleteUserForPrivate(request request.UserRequest) response.UserResponse {
	var resp response.UserResponse
	endpoint := fmt.Sprintf("%s/v1/private/user/%d", rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, request.ID)
	err := sendRequest("DELETE", endpoint, nil, &resp)
	if err != nil {
		resp.Code = "CLIENT_USER_DELETE_PRIVATE_001"
		resp.Message = err.Error()
	}
	return resp
}

func NewUserRepository(conf config.BaseConfig) UserRepository {
	return &userRepository{BaseConfig: conf}
}
