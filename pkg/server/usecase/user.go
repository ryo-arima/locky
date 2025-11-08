package usecase

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/logger"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type UserUsecase interface {
	GetUsers(c *gin.Context) ([]response.User, error)
	CreateUser(c *gin.Context, req request.UserRequest) (*response.User, error)
	UpdateUser(c *gin.Context, req request.UserRequest) (*response.User, error)
	DeleteUser(c *gin.Context, req request.UserRequest) error
	ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]response.User, error)
	CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uc *userUsecase) GetUsers(c *gin.Context) ([]response.User, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UUGU1, reqID, "Getting all users")

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

	logger.Info(code.UUGU1, reqID, "Users retrieved successfully")
	return responseUsers, nil
}

func (uc *userUsecase) CreateUser(c *gin.Context, req request.UserRequest) (*response.User, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UUCR1, reqID, "Creating user: "+req.Email)

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

	logger.Info(code.UUCR2, reqID, "User created successfully: "+createdUser.UUID)

	// Convert model to response
	return &response.User{
		ID:    createdUser.ID,
		UUID:  createdUser.UUID,
		Email: createdUser.Email,
		Name:  createdUser.Name,
	}, nil
}

func (uc *userUsecase) UpdateUser(c *gin.Context, req request.UserRequest) (*response.User, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UUUP1, reqID, "Updating user: "+req.UUID)

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

	logger.Info(code.UUUP2, reqID, "User updated successfully: "+updatedUser.UUID)

	// Convert model to response
	return &response.User{
		ID:    updatedUser.ID,
		UUID:  updatedUser.UUID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
	}, nil
}

func (uc *userUsecase) DeleteUser(c *gin.Context, req request.UserRequest) error {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UUDL1, reqID, "Deleting user: "+req.UUID)

	// Convert request to model
	user := model.Users{
		UUID: req.UUID,
	}

	// Call repository
	uc.userRepo.DeleteUser(c, user)

	logger.Info(code.UUDL1, reqID, "User deleted successfully: "+req.UUID)
	return nil
}

func (uc *userUsecase) ListUsers(c *gin.Context, filter repository.UserQueryFilter) ([]response.User, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UULS1, reqID, "Listing users with filter")

	users, err := uc.userRepo.ListUsers(c, filter)
	if err != nil {
		logger.Error(code.UULS1, reqID, "Failed to list users: "+err.Error())
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

	logger.Info(code.UULS1, reqID, "Users listed successfully")
	return responseUsers, nil
}

func (uc *userUsecase) CountUsers(c *gin.Context, filter repository.UserQueryFilter) (int64, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.UUCT1, reqID, "Counting users with filter")

	count, err := uc.userRepo.CountUsers(c, filter)
	if err != nil {
		logger.Error(code.UUCT1, reqID, "Failed to count users: "+err.Error())
		return 0, err
	}

	logger.Info(code.UUCT1, reqID, "Users counted successfully")
	return count, nil
}
