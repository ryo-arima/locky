package usecase

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
)

// Helper function to convert code.MCode to global.MCode
func toGlobalMCode(c code.MCode) global.MCode {
	return global.MCode{
		Code:    c.Code,
		Message: c.Message,
	}
}

type User interface {
	GetUsers(c *gin.Context) ([]response.User, error)
	CreateUser(c *gin.Context, req request.User) (*response.User, error)
	UpdateUser(c *gin.Context, req request.User) (*response.User, error)
	DeleteUser(c *gin.Context, req request.User) error
	ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]response.User, error)
	CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error)
	GetUserModelByEmail(c *gin.Context, email string) (*model.Users, error)
}

type user struct {
	userRepo repository.User
}

func NewUser(userRepo repository.User) User {
	return &user{
		userRepo: userRepo,
	}
}

func (uc *user) GetUsers(c *gin.Context) ([]response.User, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUGU1), "Getting all users")

	users := uc.userRepo.GetUsers(c)

	responseUsers := make([]response.User, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, response.User{
			ID:    user.ID,
			UUID:  user.UUID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	INFO(reqID, toGlobalMCode(code.UUGU1), "Users retrieved successfully")
	return responseUsers, nil
}

func (uc *user) CreateUser(c *gin.Context, req request.User) (*response.User, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUCR1), "Creating user: "+req.Email)

	// Convert request to model
	now := time.Now()
	user := model.Users{
		UUID:      req.UUID,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	// Call repository
	createdUser := uc.userRepo.CreateUser(c, user)

	INFO(reqID, toGlobalMCode(code.UUCR2), "User created successfully: "+createdUser.UUID)

	// Convert model to response
	return &response.User{
		ID:    createdUser.ID,
		UUID:  createdUser.UUID,
		Email: createdUser.Email,
		Name:  createdUser.Name,
	}, nil
}

func (uc *user) UpdateUser(c *gin.Context, req request.User) (*response.User, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUUP1), "Updating user: "+req.UUID)

	// Convert request to model
	now := time.Now()
	user := model.Users{
		UUID:      req.UUID,
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		UpdatedAt: &now,
	}

	// Call repository
	updatedUser := uc.userRepo.UpdateUser(c, user)

	INFO(reqID, toGlobalMCode(code.UUUP2), "User updated successfully: "+updatedUser.UUID)

	// Convert model to response
	return &response.User{
		ID:    updatedUser.ID,
		UUID:  updatedUser.UUID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
	}, nil
}

func (uc *user) DeleteUser(c *gin.Context, req request.User) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUDL1), "Deleting user: "+req.UUID)

	// Convert request to model
	user := model.Users{
		UUID: req.UUID,
	}

	// Call repository
	uc.userRepo.DeleteUser(c, user)

	INFO(reqID, toGlobalMCode(code.UUDL1), "User deleted successfully: "+req.UUID)
	return nil
}

func (uc *user) ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]response.User, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UULS1), "Listing users with filter")

	users, err := uc.userRepo.ListUsers(c, filter)
	if err != nil {
		ERROR(reqID, toGlobalMCode(code.UULS1), "Failed to list users: "+err.Error())
		return nil, err
	}

	responseUsers := make([]response.User, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, response.User{
			ID:    user.ID,
			UUID:  user.UUID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	INFO(reqID, toGlobalMCode(code.UULS1), "Users listed successfully")
	return responseUsers, nil
}

func (uc *user) CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUCT1), "Counting users with filter")

	count, err := uc.userRepo.CountUsers(c, filter)
	if err != nil {
		ERROR(reqID, toGlobalMCode(code.UUCT1), "Failed to count users: "+err.Error())
		return 0, err
	}

	INFO(reqID, toGlobalMCode(code.UUCT1), "Users counted successfully")
	return count, nil
}

func (uc *user) GetUserModelByEmail(c *gin.Context, email string) (*model.Users, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, toGlobalMCode(code.UUGU1), "Getting user model by email: "+email)

	users, err := uc.userRepo.ListUsers(c, repository.UserQueryFilter{Email: &email, Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return &users[0], nil
}
