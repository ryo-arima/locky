package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// UserControllerForPrivate provides admin-only user management endpoints.
//
// This interface exposes administrative operations that require
// private (admin) permissions:
//   - GetUsers: List all users
//   - CreateUser: Create a user
//   - UpdateUser: Update a user
//   - DeleteUser: Delete a user
type UserControllerForPrivate interface {
	GetUsers(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CountUsers(c *gin.Context)
}

type userControllerForPrivate struct {
	UserRepository   repository.UserRepository
	CommonRepository repository.CommonRepository
}

// GetUsers lists all users (admin only).
//
// Route: GET /v1/private/users
// Security: Bearer token (admin)
func (userController userControllerForPrivate) GetUsers(c *gin.Context) {
	filter := repository.UserQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("name"); v != "" {
		filter.Name = &v
	}
	if v := c.Query("name_prefix"); v != "" {
		filter.NamePrefix = &v
	}
	if v := c.Query("name_like"); v != "" {
		filter.NameLike = &v
	}
	if v := c.Query("email"); v != "" {
		filter.Email = &v
	}
	if v := c.Query("email_prefix"); v != "" {
		filter.EmailPrefix = &v
	}
	if v := c.Query("email_like"); v != "" {
		filter.EmailLike = &v
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Offset = n
		}
	}
	users, err := userController.UserRepository.ListUsers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_GET__FOR__002", Message: err.Error(), Users: []response.User{}})
		return
	}
	respUsers := make([]response.User, 0, len(users))
	for _, u := range users {
		respUsers = append(respUsers, response.User{ID: u.ID, UUID: u.UUID, Email: u.Email, Name: u.Name, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, DeletedAt: u.DeletedAt})
	}
	c.JSON(http.StatusOK, &response.UserResponse{Code: "SUCCESS", Message: "Users retrieved successfully", Users: respUsers})
}

// CreateUser creates a new user (admin only).
//
// Route: POST /v1/private/users
// Security: Bearer token (admin)
func (userController userControllerForPrivate) CreateUser(c *gin.Context) {
	// swagger:operation POST /private/users users createUserPrivate
	// ---
	// summary: Create a new user.
	// description: Create a new user with the provided information.
	// parameters:
	// - name: user
	//   in: body
	//   description: The user to create.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/UserRequest"
	// responses:
	//   "200":
	//     description: The created user.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__001", Message: err.Error(), Users: []response.User{}})
		return
	}
	var userModel model.Users
	// Convert model.Users to request.UserRequest
	userRequest = repository.ConvertModelToRequest(userModel)
	res := userController.UserRepository.CreateUser(userRequest)
	c.JSON(http.StatusOK, res)
}

// UpdateUser updates a user by ID (admin only).
//
// Route: PUT /v1/private/users/{id}
// Security: Bearer token (admin)
func (userController userControllerForPrivate) UpdateUser(c *gin.Context) {
	// swagger:operation PUT /private/users/{id} users updateUserPrivate
	// ---
	// summary: Update a user.
	// description: Update a user with the provided information.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the user to update.
	//   required: true
	//   type: integer
	// - name: user
	//   in: body
	//   description: The user to update.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/UserRequest"
	// responses:
	//   "200":
	//     description: The updated user.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__001", Message: err.Error(), Users: []response.User{}})
		return
	}
	var userModel model.Users
	// Convert model.Users to request.UserRequest
	userRequest = repository.ConvertModelToRequest(userModel)
	res := userController.UserRepository.UpdateUser(userRequest)
	c.JSON(http.StatusOK, res)
}

// DeleteUser deletes a user by ID (admin only).
//
// Route: DELETE /v1/private/users/{id}
// Security: Bearer token (admin)
func (userController userControllerForPrivate) DeleteUser(c *gin.Context) {
	// swagger:operation DELETE /private/users/{id} users deleteUserPrivate
	// ---
	// summary: Delete a user.
	// description: Delete a user by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the user to delete.
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: The deleted user.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__001", Message: err.Error(), Users: []response.User{}})
		return
	}
	var uuid string
	uuidRequest := repository.ConvertUUIDToRequest(uuid)
	res := userController.UserRepository.DeleteUser(uuidRequest)
	c.JSON(http.StatusOK, res)
}

// CountUsers counts the number of users matching the filter (admin only).
//
// Route: GET /v1/private/users/count
// Security: Bearer token (admin)
func (userController userControllerForPrivate) CountUsers(c *gin.Context) {
	filter := repository.UserQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("name"); v != "" {
		filter.Name = &v
	}
	if v := c.Query("name_prefix"); v != "" {
		filter.NamePrefix = &v
	}
	if v := c.Query("name_like"); v != "" {
		filter.NameLike = &v
	}
	if v := c.Query("email"); v != "" {
		filter.Email = &v
	}
	if v := c.Query("email_prefix"); v != "" {
		filter.EmailPrefix = &v
	}
	if v := c.Query("email_like"); v != "" {
		filter.EmailLike = &v
	}
	cnt, err := userController.UserRepository.CountUsers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "SERVER_CONTROLLER_COUNT__FOR__001", "message": err.Error(), "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "Count retrieved", "count": cnt})
}

// NewUserControllerForPrivate creates a new private (admin) user controller.
//
// Parameters:
//   - userRepository: User data repository
//   - commonRepository: Common services repository (e.g., auth)
//
// Returns:
//   - UserControllerForPrivate: Configured private controller instance
func NewUserControllerForPrivate(userRepository repository.UserRepository, commonRepository repository.CommonRepository) UserControllerForPrivate {
	return &userControllerForPrivate{UserRepository: userRepository, CommonRepository: commonRepository}
}
