package repository

import (
	"strings"
	"time"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
)

type UserRepository interface {
	GetUsers() []model.Users
	CreateUser(user request.UserRequest) model.Users
	UpdateUser(user request.UserRequest) model.Users
	DeleteUser(user request.UserRequest) model.Users
	ListUsers(filter UserQueryFilter) ([]model.Users, error)
	CountUsers(filter UserQueryFilter) (int64, error)
}

type userRepository struct {
	BaseConfig config.BaseConfig
}

func (userRepository userRepository) GetUsers() []model.Users {
	var users []model.Users
	userRepository.BaseConfig.DBConnection.Find(&users)
	return users
}

func (userRepository userRepository) CreateUser(user request.UserRequest) model.Users {
	now := time.Now()
	newUser := model.Users{
		UUID:      user.UUID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
		DeletedAt: nil,
	}

	if err := userRepository.BaseConfig.DBConnection.Create(&newUser).Error; err != nil {
		return model.Users{}
	}
	return newUser
}

func (userRepository userRepository) UpdateUser(user request.UserRequest) model.Users {
	var existingUser model.Users
	if err := userRepository.BaseConfig.DBConnection.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return model.Users{}
	}

	now := time.Now()
	existingUser.Password = user.Password
	existingUser.UpdatedAt = &now

	if err := userRepository.BaseConfig.DBConnection.Save(&existingUser).Error; err != nil {
		return model.Users{}
	}
	return existingUser
}

func (userRepository userRepository) DeleteUser(user request.UserRequest) model.Users {
	var existingUser model.Users
	if err := userRepository.BaseConfig.DBConnection.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		return model.Users{}
	}

	if err := userRepository.BaseConfig.DBConnection.Delete(&existingUser).Error; err != nil {
		return model.Users{}
	}
	return existingUser
}

// ConvertModelToRequest converts a model.Users object to a request.UserRequest object.
func ConvertModelToRequest(user model.Users) request.UserRequest {
	return request.UserRequest{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
	}
}

// ConvertUUIDToRequest creates a request.UserRequest object from a UUID.
func ConvertUUIDToRequest(uuid string) request.UserRequest {
	return request.UserRequest{
		Email: uuid, // Assuming Email is used as a unique identifier.
	}
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

// ListUsers filter + pagination retrieval
func (userRepository userRepository) ListUsers(filter UserQueryFilter) ([]model.Users, error) {
	filter.normalize()
	db := userRepository.BaseConfig.DBConnection
	if db == nil {
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
		return []model.Users{}, err
	}
	return users, nil
}

// CountUsers count with filter conditions
func (userRepository userRepository) CountUsers(filter UserQueryFilter) (int64, error) {
	db := userRepository.BaseConfig.DBConnection
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
		return 0, err
	}
	return cnt, nil
}

func NewUserRepository(conf config.BaseConfig) UserRepository {
	return &userRepository{BaseConfig: conf}
}
