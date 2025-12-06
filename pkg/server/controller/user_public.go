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
	"github.com/ryo-arima/locky/pkg/global"
	share "github.com/ryo-arima/locky/pkg/server/share"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

// Local logger aliases for cleaner logging code
var (
	INFO  = share.GetServerLogger().INFO
	WARN  = share.GetServerLogger().WARN
	ERROR = share.GetServerLogger().ERROR
)

// toGlobalMCode converts code.MCode to global.MCode
func toGlobalMCode(c code.MCode) global.MCode {
	return global.MCode{
		Code:    c.Code,
		Message: c.Message,
	}
}

// UserController provides public user management endpoints.
//
// This interface handles user operations that don't require prior authentication,
// primarily user registration functionality.
//
// Available endpoints:
//   - CreateUser: Handles new user registration
//   - GetUsers: Retrieves user list (public access)
type UserPublic interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

type userPublic struct {
	UserUsecase   usecase.User
	CommonUsecase usecase.Common
	conf          config.BaseConfig
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
func (rcvr userPublic) CreateUser(c *gin.Context) {
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
	var userRequest request.User
	requestID := share.GetRequestID(c)
	method := c.Request.Method
	path := c.Request.URL.Path

	// Log request received
	INFO(requestID, toGlobalMCode(code.UCPCU0), method+" "+path)

	if err := c.Bind(&userRequest); err != nil {
		WARN(requestID, toGlobalMCode(code.UCPCU1), "Bind request failed")
		res := &response.Users{Code: code.UCPCU001.Code, Message: err.Error(), Users: []response.User{}}
		WARN(requestID, toGlobalMCode(code.UCPCU6), fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate required fields
	if userRequest.Email == "" || userRequest.Name == "" || userRequest.Password == "" {
		WARN(requestID, toGlobalMCode(code.UCPCU2), "Required fields missing")
		res := &response.Users{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: "email, name, and password are required", Users: []response.User{}}
		WARN(requestID, toGlobalMCode(code.UCPCU6), fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Check for duplicate email
	users, _ := rcvr.UserUsecase.GetUsers(c)
	for _, user := range users {
		if user.Email == userRequest.Email {
			WARN(requestID, toGlobalMCode(code.UCPCU3), "Email already exists")
			res := &response.Users{Code: "SERVER_CONTROLLER_CREATE__FOR__003", Message: "email already exists", Users: []response.User{}}
			WARN(requestID, toGlobalMCode(code.UCPCU6), fmt.Sprintf("%d", http.StatusBadRequest))
			c.JSON(http.StatusBadRequest, res)
			return
		}
	}

	// Generate UUID locally
	userRequest.UUID = uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		ERROR(requestID, toGlobalMCode(code.UCPCU4), "Password hash failed")
		res := &response.Users{Code: "SERVER_CONTROLLER_CREATE__FOR__004", Message: "failed to hash password", Users: []response.User{}}
		ERROR(requestID, toGlobalMCode(code.UCPCU6), fmt.Sprintf("%d", http.StatusInternalServerError))
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	userRequest.Password = string(hashedPassword)

	createdUserPtr, _ := rcvr.UserUsecase.CreateUser(c, userRequest)

	// Convert to response format
	userResponse := response.Users{
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

	INFO(requestID, toGlobalMCode(code.UCPCU5), "User created")

	INFO(requestID, toGlobalMCode(code.UCPCU6), fmt.Sprintf("%d", http.StatusOK))
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
func (rcvr userPublic) GetUsers(c *gin.Context) {
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/UserResponse"
	var userRequest request.User
	requestID := share.GetRequestID(c)
	method := c.Request.Method
	path := c.Request.URL.Path

	// Log request received
	INFO(requestID, toGlobalMCode(code.UCPGU0), method+" "+path)

	if err := c.Bind(&userRequest); err != nil {
		WARN(requestID, toGlobalMCode(code.UCPGU1), "Bind request failed")
		res := &response.Users{Code: "SERVER_CONTROLLER_GET__FOR__001", Message: err.Error(), Users: []response.User{}}
		WARN(requestID, toGlobalMCode(code.UCPGU3), fmt.Sprintf("%d", http.StatusBadRequest))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	users, _ := rcvr.UserUsecase.GetUsers(c)

	userResponse := response.Users{
		Code:    "SUCCESS",
		Message: "Users retrieved successfully",
		Users:   users,
	}

	INFO(requestID, toGlobalMCode(code.UCPGU2), "Users retrieved")

	INFO(requestID, toGlobalMCode(code.UCPGU3), fmt.Sprintf("%d", http.StatusOK))
	c.JSON(http.StatusOK, userResponse)
}

func NewUserPublic(userUsecase usecase.User, commonUsecase usecase.Common, conf config.BaseConfig) UserPublic {
	return &userPublic{
		UserUsecase:   userUsecase,
		CommonUsecase: commonUsecase,
		conf:          conf,
	}
}
