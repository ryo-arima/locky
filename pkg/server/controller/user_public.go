package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// UserControllerForPublic provides public user management endpoints.
//
// This interface handles user operations that don't require prior authentication,
// primarily user registration functionality.
//
// Available endpoints:
//   - CreateUser: Handles new user registration
//   - GetUsers: Retrieves user list (public access)
type UserControllerForPublic interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

type userControllerForPublic struct {
	UserRepository   repository.UserRepository
	CommonRepository repository.CommonRepository
}

// CreateUser handles new user registration.
//
// This endpoint allows anonymous users to register new accounts.
// It validates the input, checks for duplicate emails, generates UUIDs,
// hashes passwords using bcrypt, and creates the user record.
//
// Route: POST /v1/public/users
// Security: No authentication required
//
// swagger:route POST /public/users Public Users createUser
//
// # Create a new user account
//
// Allows anonymous users to register new accounts with email, name, and password.
//
// Responses:
//
//	200: userResponse
//	400: errorResponse
//	500: errorResponse
func (userController userControllerForPublic) CreateUser(c *gin.Context) {
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
	//   "500":
	//     description: Internal server error.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	if err := c.Bind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__001", Message: err.Error(), Users: []response.User{}})
		return
	}

	// Validate required fields
	if userRequest.Email == "" || userRequest.Name == "" || userRequest.Password == "" {
		c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: "email, name, and password are required", Users: []response.User{}})
		return
	}

	// Check for duplicate email
	users := userController.UserRepository.GetUsers()
	for _, user := range users {
		if user.Email == userRequest.Email {
			c.JSON(http.StatusBadRequest, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__003", Message: "email already exists", Users: []response.User{}})
			return
		}
	}

	// Generate UUID locally
	userRequest.UUID = uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__004", Message: "failed to hash password", Users: []response.User{}})
		return
	}
	userRequest.Password = string(hashedPassword)

	createdUser := userController.UserRepository.CreateUser(userRequest)

	// Convert to response format
	userResponse := response.UserResponse{
		Code:    "SUCCESS",
		Message: "User created successfully",
		Users: []response.User{
			{
				ID:    createdUser.ID,
				UUID:  createdUser.UUID,
				Email: createdUser.Email,
				Name:  createdUser.Name,
			},
		},
	}

	c.JSON(http.StatusOK, userResponse)
}

// GetUsers retrieves a list of all users.
//
// This endpoint returns a list of all users in the system.
// Note: This is a public endpoint and doesn't require authentication.
//
// Route: GET /v1/public/users
// Security: No authentication required
//
// swagger:route GET /public/users Public Users getUsers
//
// # Get list of all users
//
// Returns a list of all users in the system (public access).
//
// Responses:
//
//	200: userResponse
//	400: errorResponse
func (userController userControllerForPublic) GetUsers(c *gin.Context) {
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

	users := userController.UserRepository.GetUsers()

	// Convert to response format
	var responseUsers []response.User
	for _, user := range users {
		responseUsers = append(responseUsers, response.User{
			ID:    user.ID,
			UUID:  user.UUID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	userResponse := response.UserResponse{
		Code:    "SUCCESS",
		Message: "Users retrieved successfully",
		Users:   responseUsers,
	}

	c.JSON(http.StatusOK, userResponse)
}

func NewUserControllerForPublic(userRepository repository.UserRepository, commonRepository repository.CommonRepository) UserControllerForPublic {
	return &userControllerForPublic{
		UserRepository:   userRepository,
		CommonRepository: commonRepository,
	}
}
