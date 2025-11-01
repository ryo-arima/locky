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

// UserControllerForInternal provides authenticated (non-admin) user operations.
//
// This interface exposes standard user operations requiring authentication:
//   - GetUsers: List users
//   - UpdateUser: Update user
//   - DeleteUser: Delete user
//   - CreateUser: Create user
type UserControllerForInternal interface {
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CreateUser(c *gin.Context)
	CountUsers(c *gin.Context) // Added: count
}

type userControllerForInternal struct {
	UserRepository   repository.UserRepository
	CommonRepository repository.CommonRepository
}

// GetUsers lists users (authenticated).
//
// Route: GET /v1/internal/users
// Security: Bearer token
func (userController userControllerForInternal) GetUsers(c *gin.Context) {
	// swagger:operation GET /internal/users users getUsersInternal
	// ---
	// summary: Get a list of users.
	// description: Get a list of all users in the system.
	// responses:
	//   "200":
	//     description: A list of users.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_GET__FOR__001", Message: err.Error(), Users: []response.User{}})
		return
	}
	// Search with query parameters (no body required)
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

// UpdateUser updates a user (authenticated).
//
// Route: PUT /v1/internal/users/{id}
// Security: Bearer token
func (userController userControllerForInternal) UpdateUser(c *gin.Context) {
	// swagger:operation PUT /internal/users/{id} users updateUserInternal
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

// DeleteUser deletes a user (authenticated).
//
// Route: DELETE /v1/internal/users/{id}
// Security: Bearer token
func (userController userControllerForInternal) DeleteUser(c *gin.Context) {
	// swagger:operation DELETE /internal/users/{id} users deleteUserInternal
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

// CreateUser creates a user (authenticated).
//
// Route: POST /v1/internal/users
// Security: Bearer token
func (userController userControllerForInternal) CreateUser(c *gin.Context) {
	// swagger:operation POST /internal/users users createUserInternal
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
	userRequest = repository.ConvertModelToRequest(userModel)
	res := userController.UserRepository.CreateUser(userRequest)
	c.JSON(http.StatusOK, res)
}

// CountUsers counts users (authenticated).
//
// Route: GET /v1/internal/users/count
// Security: Bearer token
func (userController userControllerForInternal) CountUsers(c *gin.Context) {
	// swagger:operation GET /internal/users/count users countUsersInternal
	// ---
	// summary: Get the count of users.
	// description: Get the total number of users in the system.
	// responses:
	//   "200":
	//     description: The count of users.
	//     schema:
	//       type: object
	//       properties:
	//         count:
	//           type: integer
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
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

// NewUserControllerForInternal creates a new internal user controller.
//
// Parameters:
//   - userRepository: User data repository
//   - commonRepository: Common services repository (e.g., auth)
//
// Returns:
//   - UserControllerForInternal: Configured internal controller instance
func NewUserControllerForInternal(userRepository repository.UserRepository, commonRepository repository.CommonRepository) UserControllerForInternal {
	return &userControllerForInternal{UserRepository: userRepository, CommonRepository: commonRepository}
}
