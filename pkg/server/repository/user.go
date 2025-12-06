package repository

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/share"
)

type User interface {
	GetUsers(c *gin.Context) []model.Users
	GetUserByEmail(c *gin.Context, email string) (*model.Users, error)
	CreateUser(c *gin.Context, user model.Users) model.Users
	UpdateUser(c *gin.Context, user model.Users) model.Users
	DeleteUser(c *gin.Context, user model.Users) model.Users
	ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error)
	CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error)
}

type user struct {
	BaseConfig config.BaseConfig
}

func (rcvr user) GetUsers(c *gin.Context) []model.Users {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Getting all users from database")

	var users []model.Users
	rcvr.BaseConfig.DBConnection.Find(&users)

	INFO(reqID, global.SRNRSR1, "Retrieved users from database")
	return users
}

func (rcvr user) GetUserByEmail(c *gin.Context, email string) (*model.Users, error) {
	var user model.Users
	if err := rcvr.BaseConfig.DBConnection.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (rcvr user) CreateUser(c *gin.Context, user model.Users) model.Users {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Creating user in database: "+user.Email)

	if err := rcvr.BaseConfig.DBConnection.Create(&user).Error; err != nil {
		ERROR(reqID, global.SRNRSR2, "Failed to create user: "+err.Error())
		return model.Users{}
	}

	INFO(reqID, global.SRNRSR1, "User created in database: "+user.UUID)
	return user
}

func (rcvr user) UpdateUser(c *gin.Context, user model.Users) model.Users {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Updating user in database: "+user.UUID)

	if err := rcvr.BaseConfig.DBConnection.Save(&user).Error; err != nil {
		ERROR(reqID, global.SRNRSR2, "Failed to update user: "+err.Error())
		return model.Users{}
	}

	INFO(reqID, global.SRNRSR1, "User updated in database: "+user.UUID)
	return user
}

func (rcvr user) DeleteUser(c *gin.Context, user model.Users) model.Users {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Deleting user from database: "+user.UUID)

	if err := rcvr.BaseConfig.DBConnection.Delete(&user).Error; err != nil {
		ERROR(reqID, global.SRNRSR2, "Failed to delete user: "+err.Error())
		return model.Users{}
	}

	INFO(reqID, global.SRNRSR1, "User deleted from database: "+user.UUID)
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
func (rcvr user) ListUsers(c *gin.Context, filter UserQueryFilter) ([]model.Users, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Listing users from database with filter")

	filter.normalize()
	db := rcvr.BaseConfig.DBConnection
	if db == nil {
		WARN(reqID, global.SRNRSR1, "Database connection is nil")
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
		ERROR(reqID, global.SRNRSR2, "Failed to list users: "+err.Error())
		return []model.Users{}, err
	}

	INFO(reqID, global.SRNRSR1, "Users listed from database successfully")
	return users, nil
}

// CountUsers counts users with filter conditions
func (rcvr user) CountUsers(c *gin.Context, filter UserQueryFilter) (int64, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, global.SRNRSR1, "Counting users in database with filter")

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
		ERROR(reqID, global.SRNRSR2, "Failed to count users: "+err.Error())
		return 0, err
	}

	INFO(reqID, global.SRNRSR1, "Users counted in database successfully")
	return cnt, nil
}

func NewUser(conf config.BaseConfig) User {
	return &user{BaseConfig: conf}
}
