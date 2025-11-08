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

type MemberRepository interface {
	GetMembers(c *gin.Context) []model.Members
	CreateMember(c *gin.Context, member *model.Members) *gorm.DB
	UpdateMember(c *gin.Context, member *model.Members) *gorm.DB
	DeleteMember(c *gin.Context, uuid string) *gorm.DB
	GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error)
	ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error)
	CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error)
}

type memberRepository struct {
	BaseConfig config.BaseConfig
}

func (rcvr memberRepository) GetMembers(c *gin.Context) []model.Members {
	var members []model.Members
	rcvr.BaseConfig.DBConnection.Find(&members)
	return members
}

func (rcvr memberRepository) CreateMember(c *gin.Context, member *model.Members) *gorm.DB {
	if member == nil {
		return &gorm.DB{Error: errors.New("member is nil")}
	}
	return rcvr.BaseConfig.DBConnection.Create(member)
}

func (rcvr memberRepository) UpdateMember(c *gin.Context, member *model.Members) *gorm.DB {
	if member == nil {
		return &gorm.DB{Error: errors.New("member is nil")}
	}
	return rcvr.BaseConfig.DBConnection.Model(&model.Members{}).Where("id = ?", member.ID).Updates(member)
}

func (rcvr memberRepository) DeleteMember(c *gin.Context, uuid string) *gorm.DB {
	return rcvr.BaseConfig.DBConnection.Model(&model.Members{}).Where("uuid = ?", uuid).Update("deleted_at", time.Now())
}

func (rcvr memberRepository) GetMemberByUUID(c *gin.Context, uuid string) (model.Members, error) {
	var m model.Members
	res := rcvr.BaseConfig.DBConnection.Where("uuid = ?", uuid).First(&m)
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
func (rcvr memberRepository) ListMembers(c *gin.Context, filter MemberQueryFilter) ([]model.Members, error) {
	filter.normalize()
	q := rcvr.BaseConfig.DBConnection.Model(&model.Members{})
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
func (rcvr memberRepository) CountMembers(c *gin.Context, filter MemberQueryFilter) (int64, error) {
	q := rcvr.BaseConfig.DBConnection.Model(&model.Members{})
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

func NewMemberRepository(conf config.BaseConfig) MemberRepository {
	return &memberRepository{BaseConfig: conf}
}
