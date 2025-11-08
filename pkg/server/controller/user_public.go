package controller

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/logger"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
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
	UserUsecase      usecase.UserUsecase
	CommonRepository repository.CommonRepository
	conf             config.BaseConfig
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
func (rcvr userControllerForPublic) CreateUser(c *gin.Context) {
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
	requestID := middleware.GetRequestID(c)
	method := c.Request.Method
	path := c.Request.URL.Path

	// Log request received
	logger.Info(code.UCPCU0, requestID, method+" "+path)

	if err := c.Bind(&userRequest); err != nil {
		logger.Warn(code.UCPCU1, requestID, "Bind request failed")
		res := &response.UserResponse{Code: code.UCPCU001.Code, Message: err.Error(), Users: []response.User{}}
		logger.Warn(code.UCPCU6, requestID, fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate required fields
	if userRequest.Email == "" || userRequest.Name == "" || userRequest.Password == "" {
		logger.Warn(code.UCPCU2, requestID, "Required fields missing")
		res := &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: "email, name, and password are required", Users: []response.User{}}
		logger.Warn(code.UCPCU6, requestID, fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Check for duplicate email
	users, _ := rcvr.UserUsecase.GetUsers(c)
	for _, user := range users {
		if user.Email == userRequest.Email {
			logger.Warn(code.UCPCU3, requestID, "Email already exists")
			res := &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__003", Message: "email already exists", Users: []response.User{}}
			logger.Warn(code.UCPCU6, requestID, fmt.Sprintf("%d", http.StatusBadRequest))
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	// Generate UUID locally
	userRequest.UUID = uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(code.UCPCU4, requestID, "Password hash failed")
		res := &response.UserResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__004", Message: "failed to hash password", Users: []response.User{}}
		logger.Error(code.UCPCU6, requestID, fmt.Sprintf("%d", http.StatusInternalServerError))
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	userRequest.Password = string(hashedPassword)

	createdUserPtr, _ := rcvr.UserUsecase.CreateUser(c, userRequest)

	// Convert to response format
	userResponse := response.UserResponse{
		Code:    "SUCCESS",
		Message: "User created successfully",
		Users: []response.User{
			{
				ID:    createdUserPtr.ID,
				UUID:  createdUserPtr.UUID,
				Email: createdUserPtr.Email,
				Name:  createdUserPtr.Name,
			},
		},
	}

	logger.Info(code.UCPCU5, requestID, "User created")

	logger.Info(code.UCPCU6, requestID, fmt.Sprintf("%d", http.StatusOK))
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
func (rcvr userControllerForPublic) GetUsers(c *gin.Context) {
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.UserRequest
	requestID := middleware.GetRequestID(c)
	method := c.Request.Method
	path := c.Request.URL.Path

	// Log request received
	logger.Info(code.UCPGU0, requestID, method+" "+path)

	if err := c.Bind(&userRequest); err != nil {
		logger.Warn(code.UCPGU1, requestID, "Bind request failed")
		res := &response.UserResponse{Code: "SERVER_CONTROLLER_GET__FOR__001", Message: err.Error(), Users: []response.User{}}
		logger.Warn(code.UCPGU3, requestID, fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	users, _ := rcvr.UserUsecase.GetUsers(c)

	userResponse := response.UserResponse{
		Code:    "SUCCESS",
		Message: "Users retrieved successfully",
		Users:   users,
	}

	logger.Info(code.UCPGU2, requestID, "Users retrieved")

	logger.Info(code.UCPGU3, requestID, fmt.Sprintf("%d", http.StatusOK))
	c.JSON(http.StatusOK, userResponse)
}

func NewUserControllerForPublic(userUsecase usecase.UserUsecase, commonRepository repository.CommonRepository, conf config.BaseConfig) UserControllerForPublic {
	return &userControllerForPublic{
		UserUsecase:      userUsecase,
		CommonRepository: commonRepository,
		conf:             conf,
	}
}
