package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"gorm.io/gorm"
)

type Member interface {
	GetMembers(c *gin.Context) ([]response.Member, error)
	GetMemberByUUID(c *gin.Context, uuid string) (*response.Member, error)
	CreateMember(c *gin.Context, member *model.Members) (*response.Member, *gorm.DB)
	UpdateMember(c *gin.Context, member *model.Members) (*response.Member, *gorm.DB)
	DeleteMember(c *gin.Context, uuid string) *gorm.DB
	ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]response.Member, error)
	CountMembers(c *gin.Context, filter repository.MemberQueryFilter) (int64, error)
}

type member struct {
	memberRepo repository.Member
}

func NewMember(memberRepo repository.Member) Member {
	return &member{
		memberRepo: memberRepo,
	}
}

func (uc *member) GetMembers(c *gin.Context) ([]response.Member, error) {
	members := uc.memberRepo.GetMembers(c)

	var responseMembers []response.Member
	for _, member := range members {
		responseMembers = append(responseMembers, response.Member{
			ID:        member.ID,
			UUID:      member.UUID,
			GroupUUID: member.GroupUUID,
			UserUUID:  member.UserUUID,
			Role:      member.Role,
		})
	}

	return responseMembers, nil
}

func (uc *member) GetMemberByUUID(c *gin.Context, uuid string) (*response.Member, error) {
	member, err := uc.memberRepo.GetMemberByUUID(c, uuid)
	if err != nil {
		return nil, err
	}

	return &response.Member{
		ID:        member.ID,
		UUID:      member.UUID,
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
	}, nil
}

func (uc *member) CreateMember(c *gin.Context, member *model.Members) (*response.Member, *gorm.DB) {
	resDB := uc.memberRepo.CreateMember(c, member)
	if resDB.Error != nil {
		return nil, resDB
	}

	return &response.Member{
		ID:        member.ID,
		UUID:      member.UUID,
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
	}, resDB
}

func (uc *member) UpdateMember(c *gin.Context, member *model.Members) (*response.Member, *gorm.DB) {
	resDB := uc.memberRepo.UpdateMember(c, member)
	if resDB.Error != nil {
		return nil, resDB
	}

	return &response.Member{
		ID:        member.ID,
		UUID:      member.UUID,
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
	}, resDB
}

func (uc *member) DeleteMember(c *gin.Context, uuid string) *gorm.DB {
	return uc.memberRepo.DeleteMember(c, uuid)
}

func (uc *member) ListMembers(c *gin.Context, filter repository.MemberQueryFilter) ([]response.Member, error) {
	members, err := uc.memberRepo.ListMembers(c, filter)
	if err != nil {
		return nil, err
	}

	var responseMembers []response.Member
	for _, member := range members {
		responseMembers = append(responseMembers, response.Member{
			ID:        member.ID,
			UUID:      member.UUID,
			GroupUUID: member.GroupUUID,
			UserUUID:  member.UserUUID,
			Role:      member.Role,
		})
	}

	return responseMembers, nil
}

func (uc *member) CountMembers(c *gin.Context, filter repository.MemberQueryFilter) (int64, error) {
	return uc.memberRepo.CountMembers(c, filter)
}
