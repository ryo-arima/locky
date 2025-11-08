package repository

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/code"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/logger"
)

type UserRepository interface {
	GetUsers(c *gin.Context) []model.Users
	CreateUser(c *gin.Context, user model.Users) model.Users
	UpdateUser(c *gin.Context, user model.Users) model.Users
	DeleteUser(c *gin.Context, user model.Users) model.Users
	ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error)
	CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error)
}

type userRepository struct {
	BaseConfig config.BaseConfig
}

func (rcvr userRepository) GetUsers(c *gin.Context) []model.Users {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RURP1, reqID, "Getting all users from database")

	var users []model.Users
	rcvr.BaseConfig.DBConnection.Find(&users)

	logger.Info(code.RURP1, reqID, "Retrieved users from database")
	return users
}

func (rcvr userRepository) CreateUser(c *gin.Context, user model.Users) model.Users {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RUCR1, reqID, "Creating user in database: "+user.Email)

	if err := rcvr.BaseConfig.DBConnection.Create(&user).Error; err != nil {
		logger.Error(code.RUCR1, reqID, "Failed to create user: "+err.Error())
		return model.Users{}
	}

	logger.Info(code.RUCR1, reqID, "User created in database: "+user.UUID)
	return user
}

func (rcvr userRepository) UpdateUser(c *gin.Context, user model.Users) model.Users {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RUUP1, reqID, "Updating user in database: "+user.UUID)

	if err := rcvr.BaseConfig.DBConnection.Save(&user).Error; err != nil {
		logger.Error(code.RUUP1, reqID, "Failed to update user: "+err.Error())
		return model.Users{}
	}

	logger.Info(code.RUUP1, reqID, "User updated in database: "+user.UUID)
	return user
}

func (rcvr userRepository) DeleteUser(c *gin.Context, user model.Users) model.Users {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RUDL1, reqID, "Deleting user from database: "+user.UUID)

	if err := rcvr.BaseConfig.DBConnection.Delete(&user).Error; err != nil {
		logger.Error(code.RUDL1, reqID, "Failed to delete user: "+err.Error())
		return model.Users{}
	}

	logger.Info(code.RUDL1, reqID, "User deleted from database: "+user.UUID)
	return user
}

// UserQueryFilter: search conditions for GET /users
type UserQueryFilter struct {
	ID          *uint
	UUID        *string
	Name        *string // exact match
	NamePrefix  *string
	NameLike    *string
	Email       *string // exact match
	EmailPrefix *string
	EmailLike   *string
	Limit       int
	Offset      int
}

func (f *UserQueryFilter) normalize() {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}

// ListUsers retrieves users with filter and pagination
func (rcvr userRepository) ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RULS1, reqID, "Listing users from database with filter")

	filter.normalize()
	db := rcvr.BaseConfig.DBConnection
	if db == nil {
		logger.Warn(code.RULS1, reqID, "Database connection is nil")
		return []model.Users{}, nil
	}
	q := db.Model(&model.Users{})
	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.Name != nil {
		q = q.Where("name = ?", *filter.Name)
	}
	if filter.NamePrefix != nil {
		q = q.Where("name LIKE ?", strings.TrimRight(*filter.NamePrefix, "%")+"%")
	}
	if filter.NameLike != nil {
		q = q.Where("name LIKE ?", "%"+*filter.NameLike+"%")
	}
	if filter.Email != nil {
		q = q.Where("email = ?", *filter.Email)
	}
	if filter.EmailPrefix != nil {
		q = q.Where("email LIKE ?", strings.TrimRight(*filter.EmailPrefix, "%")+"%")
	}
	if filter.EmailLike != nil {
		q = q.Where("email LIKE ?", "%"+*filter.EmailLike+"%")
	}
	q = q.Limit(filter.Limit).Offset(filter.Offset)
	var users []model.Users
	if err := q.Find(&users).Error; err != nil {
		logger.Error(code.RULS1, reqID, "Failed to list users: "+err.Error())
		return []model.Users{}, err
	}

	logger.Info(code.RULS1, reqID, "Users listed from database successfully")
	return users, nil
}

// CountUsers counts users with filter conditions
func (rcvr userRepository) CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error) {
	requestID, _ := c.Get("requestID")
	reqID := requestID.(string)
	logger.Info(code.RUCT1, reqID, "Counting users in database with filter")

	db := rcvr.BaseConfig.DBConnection
	if db == nil {
		return 0, nil
	}
	q := db.Model(&model.Users{})
	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.Name != nil {
		q = q.Where("name = ?", *filter.Name)
	}
	if filter.NamePrefix != nil {
		q = q.Where("name LIKE ?", strings.TrimRight(*filter.NamePrefix, "%")+"%")
	}
	if filter.NameLike != nil {
		q = q.Where("name LIKE ?", "%"+*filter.NameLike+"%")
	}
	if filter.Email != nil {
		q = q.Where("email = ?", *filter.Email)
	}
	if filter.EmailPrefix != nil {
		q = q.Where("email LIKE ?", strings.TrimRight(*filter.EmailPrefix, "%")+"%")
	}
	if filter.EmailLike != nil {
		q = q.Where("email LIKE ?", "%"+*filter.EmailLike+"%")
	}
	var cnt int64
	if err := q.Count(&cnt).Error; err != nil {
		logger.Error(code.RUCT1, reqID, "Failed to count users: "+err.Error())
		return 0, err
	}

	logger.Info(code.RUCT1, reqID, "Users counted in database successfully")
	return cnt, nil
}

func NewUserRepository(conf config.BaseConfig) UserRepository {
	return &userRepository{BaseConfig: conf}
}
