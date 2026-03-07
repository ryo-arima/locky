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

type Member interface {
	GetMembers(c *gin.Context) []model.Members
	CreateMember(c *gin.Context, member *model.Members) error
	UpdateMember(c *gin.Context, member *model.Members) error
	DeleteMember(c *gin.Context, uuid string) error
	GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error)
	ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error)
	CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error)
	// WithTx returns a new Member repository that uses the given transaction.
	// Use this to participate in a transaction managed by the usecase layer.
	WithTx(tx *gorm.DB) Member
}

type member struct {
	db *gorm.DB
}

// WithTx returns a new Member backed by the given transaction.
func (rcvr *member) WithTx(tx *gorm.DB) Member {
	return &member{db: tx}
}

func (rcvr *member) GetMembers(c *gin.Context) []model.Members {
	var members []model.Members
	rcvr.db.Find(&members)
	return members
}

func (rcvr *member) CreateMember(c *gin.Context, member *model.Members) error {
	if member == nil {
		return errors.New("member is nil")
	}
	return rcvr.db.Create(member).Error
}

func (rcvr *member) UpdateMember(c *gin.Context, member *model.Members) error {
	if member == nil {
		return errors.New("member is nil")
	}
	return rcvr.db.Model(&model.Members{}).Where("id = ?", member.ID).Updates(member).Error
}

func (rcvr *member) DeleteMember(c *gin.Context, uuid string) error {
	return rcvr.db.Model(&model.Members{}).Where("uuid = ?", uuid).Update("deleted_at", time.Now()).Error
}

func (rcvr *member) GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error) {
	var m model.Members
	res := rcvr.db.Where("uuid = ?", uuid).First(&m)
	if res.Error != nil {
		return model.Members{}, res.Error
	}
	return m, nil
}

// MemberQueryFilter: member search/pagination conditions
type MemberQueryFilter struct {
	ID         *uint
	UUID       *string
	GroupUUID  *string
	UserUUID   *string
	Role       *string
	RolePrefix *string
	RoleLike   *string
	Limit      int
	Offset     int
}

func (f *MemberQueryFilter) normalize() {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}

// ListMembers filter + pagination
func (rcvr *member) ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error) {
	filter.normalize()
	q := rcvr.db.Model(&model.Members{})
	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.GroupUUID != nil {
		q = q.Where("group_uuid = ?", *filter.GroupUUID)
	}
	if filter.UserUUID != nil {
		q = q.Where("user_uuid = ?", *filter.UserUUID)
	}
	if filter.Role != nil {
		q = q.Where("role = ?", *filter.Role)
	}
	if filter.RolePrefix != nil {
		q = q.Where("role LIKE ?", strings.TrimRight(*filter.RolePrefix, "%")+"%")
	}
	if filter.RoleLike != nil {
		q = q.Where("role LIKE ?", "%"+*filter.RoleLike+"%")
	}
	q = q.Limit(filter.Limit).Offset(filter.Offset)
	var list []model.Members
	if err := q.Find(&list).Error; err != nil {
		return []model.Members{}, err
	}
	return list, nil
}

// CountMembers get count
func (rcvr *member) CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error) {
	q := rcvr.db.Model(&model.Members{})
	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.GroupUUID != nil {
		q = q.Where("group_uuid = ?", *filter.GroupUUID)
	}
	if filter.UserUUID != nil {
		q = q.Where("user_uuid = ?", *filter.UserUUID)
	}
	if filter.Role != nil {
		q = q.Where("role = ?", *filter.Role)
	}
	if filter.RolePrefix != nil {
		q = q.Where("role LIKE ?", strings.TrimRight(*filter.RolePrefix, "%")+"%")
	}
	if filter.RoleLike != nil {
		q = q.Where("role LIKE ?", "%"+*filter.RoleLike+"%")
	}
	var cnt int64
	if err := q.Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func NewMember(conf config.BaseConfig) Member {
	return &member{db: conf.DBConnection}
}
