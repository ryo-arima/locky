package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
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
	UserUsecase      usecase.UserUsecase
	CommonRepository repository.CommonRepository
}

// GetUsers lists users (authenticated).
//
// Route: GET /v1/internal/users
// Security: Bearer token
func (rcvr userControllerForInternal) GetUsers(c *gin.Context) {
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
	users, err := rcvr.UserUsecase.ListUsers(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_GET__FOR__002", Message: err.Error(), Users: []response.User{}})
		return
	}
	c.JSON(http.StatusOK, &response.UserResponse{Code: "SUCCESS", Message: "Users retrieved successfully", Users: users})
}

// UpdateUser updates a user (authenticated).
//
// Route: PUT /v1/internal/users/{id}
// Security: Bearer token
func (rcvr userControllerForInternal) UpdateUser(c *gin.Context) {
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

	updatedUser, err := rcvr.UserUsecase.UpdateUser(c, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__002", Message: err.Error(), Users: []response.User{}})
		return
	}
	c.JSON(http.StatusOK, &response.UserResponse{Code: "SUCCESS", Message: "User updated successfully", Users: []response.User{*updatedUser}})
}

// DeleteUser deletes a user (authenticated).
//
// Route: DELETE /v1/internal/users/{id}
// Security: Bearer token
func (rcvr userControllerForInternal) DeleteUser(c *gin.Context) {
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

	err := rcvr.UserUsecase.DeleteUser(c, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__002", Message: err.Error(), Users: []response.User{}})
		return
	}
	c.JSON(http.StatusOK, &response.UserResponse{Code: "SUCCESS", Message: "User deleted successfully", Users: []response.User{}})
}

// CreateUser creates a user (authenticated).
//
// Route: POST /v1/internal/users
// Security: Bearer token
func (rcvr userControllerForInternal) CreateUser(c *gin.Context) {
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

	createdUser, err := rcvr.UserUsecase.CreateUser(c, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: err.Error(), Users: []response.User{}})
		return
	}
	c.JSON(http.StatusOK, &response.UserResponse{Code: "SUCCESS", Message: "User created successfully", Users: []response.User{*createdUser}})
}

// CountUsers counts users (authenticated).
//
// Route: GET /v1/internal/users/count
// Security: Bearer token
func (rcvr userControllerForInternal) CountUsers(c *gin.Context) {
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
	cnt, err := rcvr.UserUsecase.CountUsers(c, filter)
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
func NewUserControllerForInternal(userUsecase usecase.UserUsecase, commonRepository repository.CommonRepository) UserControllerForInternal {
	return &userControllerForInternal{UserUsecase: userUsecase, CommonRepository: commonRepository}
}
