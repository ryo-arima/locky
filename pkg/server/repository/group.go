package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"gorm.io/gorm"
)

type GroupRepository interface {
	GetGroups() []model.Groups
	GetGroupByUUID(uuid string) (model.Groups, error)
	GetGroupByID(id uint) (model.Groups, error)
	CreateGroup(group *model.Groups) *gorm.DB
	UpdateGroup(group *model.Groups) *gorm.DB
	DeleteGroup(uuid string) *gorm.DB
	ListGroups(filter GroupQueryFilter) ([]model.Groups, error)
	CountGroups(filter GroupQueryFilter) (int64, error)
}

type groupRepository struct {
	BaseConfig config.BaseConfig
}

func (r groupRepository) GetGroups() []model.Groups {
	var groups []model.Groups
	r.BaseConfig.DBConnection.Find(&groups)
	return groups
}

func (r groupRepository) GetGroupByUUID(uuid string) (model.Groups, error) {
	var g model.Groups
	res := r.BaseConfig.DBConnection.Where("uuid = ?", uuid).First(&g)
	if res.Error != nil {
		return model.Groups{}, res.Error
	}
	return g, nil
}

func (r groupRepository) GetGroupByID(id uint) (model.Groups, error) {
	var g model.Groups
	res := r.BaseConfig.DBConnection.First(&g, id)
	if res.Error != nil {
		return model.Groups{}, res.Error
	}
	return g, nil
}

func (r groupRepository) CreateGroup(group *model.Groups) *gorm.DB {
	if group == nil {
		return &gorm.DB{Error: errors.New("group is nil")}
	}
	return r.BaseConfig.DBConnection.Create(group)
}

func (r groupRepository) UpdateGroup(group *model.Groups) *gorm.DB {
	if group == nil {
		return &gorm.DB{Error: errors.New("group is nil")}
	}
	return r.BaseConfig.DBConnection.Model(&model.Groups{}).Where("id = ?", group.ID).Updates(group)
}

func (r groupRepository) DeleteGroup(uuid string) *gorm.DB {
	return r.BaseConfig.DBConnection.Model(&model.Groups{}).Where("uuid = ?", uuid).Update("deleted_at", time.Now())
}

// GroupQueryFilter: group search/pagination conditions
type GroupQueryFilter struct {
	ID         *uint
	UUID       *string
	Name       *string
	NamePrefix *string
	NameLike   *string
	Limit      int
	Offset     int
}

func (f *GroupQueryFilter) normalize() {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}

func (r groupRepository) ListGroups(filter GroupQueryFilter) ([]model.Groups, error) {
	filter.normalize()
	q := r.BaseConfig.DBConnection.Model(&model.Groups{})
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
	q = q.Limit(filter.Limit).Offset(filter.Offset)
	var list []model.Groups
	if err := q.Find(&list).Error; err != nil {
		return []model.Groups{}, err
	}
	return list, nil
}

func (r groupRepository) CountGroups(filter GroupQueryFilter) (int64, error) {
	q := r.BaseConfig.DBConnection.Model(&model.Groups{})
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
	var cnt int64
	if err := q.Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func NewGroupRepository(conf config.BaseConfig) GroupRepository {
	return &groupRepository{BaseConfig: conf}
}
