package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"gorm.io/gorm"
)

type Group interface {
	GetGroups(c *gin.Context) []model.Groups
	GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error)
	GetGroupByID(c *gin.Context, id uint) (model.Groups, error)
	CreateGroup(c *gin.Context, group *model.Groups) error
	UpdateGroup(c *gin.Context, group *model.Groups) error
	DeleteGroup(c *gin.Context, uuid string) error
	ListGroups(c *gin.Context, filter GroupQueryFilter) ([]model.Groups, error)
	CountGroups(c *gin.Context, filter GroupQueryFilter) (int64, error)
	// WithTx returns a new Group repository that uses the given transaction.
	// Use this to participate in a transaction managed by the usecase layer.
	WithTx(tx *gorm.DB) Group
}

type group struct {
	db *gorm.DB
}

// WithTx returns a new Group backed by the given transaction.
func (rcvr *group) WithTx(tx *gorm.DB) Group {
	return &group{db: tx}
}

func (rcvr *group) GetGroups(c *gin.Context) []model.Groups {
	var groups []model.Groups
	rcvr.db.Find(&groups)
	return groups
}

func (rcvr *group) GetGroupByUUID(c *gin.Context, uuid string) (model.Groups, error) {
	var g model.Groups
	res := rcvr.db.Where("uuid = ?", uuid).First(&g)
	if res.Error != nil {
		return model.Groups{}, res.Error
	}
	return g, nil
}

func (rcvr *group) GetGroupByID(c *gin.Context, id uint) (model.Groups, error) {
	var g model.Groups
	res := rcvr.db.First(&g, id)
	if res.Error != nil {
		return model.Groups{}, res.Error
	}
	return g, nil
}

func (rcvr *group) CreateGroup(c *gin.Context, group *model.Groups) error {
	if group == nil {
		return errors.New("group is nil")
	}
	return rcvr.db.Create(group).Error
}

func (rcvr *group) UpdateGroup(c *gin.Context, group *model.Groups) error {
	if group == nil {
		return errors.New("group is nil")
	}
	return rcvr.db.Model(&model.Groups{}).Where("id = ?", group.ID).Updates(group).Error
}

func (rcvr *group) DeleteGroup(c *gin.Context, uuid string) error {
	return rcvr.db.Model(&model.Groups{}).Where("uuid = ?", uuid).Update("deleted_at", time.Now()).Error
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

func (rcvr *group) ListGroups(c *gin.Context, filter GroupQueryFilter) ([]model.Groups, error) {
	filter.normalize()
	q := rcvr.db.Model(&model.Groups{})
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

func (rcvr *group) CountGroups(c *gin.Context, filter GroupQueryFilter) (int64, error) {
	q := rcvr.db.Model(&model.Groups{})
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

func NewGroup(conf config.BaseConfig) Group {
	return &group{db: conf.DBConnection}
}
