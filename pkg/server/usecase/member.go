package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
)

type Member interface {
	GetMembers(c *gin.Context) ([]response.Member, error)
	GetMemberByUUID(c *gin.Context, uuid string) (*response.Member, error)
	CreateMember(c *gin.Context, member *model.Members) (*response.Member, error)
	UpdateMember(c *gin.Context, member *model.Members) (*response.Member, error)
	DeleteMember(c *gin.Context, uuid string) error
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

func (uc *member) CreateMember(c *gin.Context, member *model.Members) (*response.Member, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CreateMember called")
	if err := uc.memberRepo.CreateMember(c, member); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to create member: "+err.Error())
		return nil, err
	}

	INFO(reqID, Mcode(SRNRSR1), "CreateMember succeeded")
	return &response.Member{
		ID:        member.ID,
		UUID:      member.UUID,
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
	}, nil
}

func (uc *member) UpdateMember(c *gin.Context, member *model.Members) (*response.Member, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "UpdateMember called")
	if err := uc.memberRepo.UpdateMember(c, member); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to update member: "+err.Error())
		return nil, err
	}

	INFO(reqID, Mcode(SRNRSR1), "UpdateMember succeeded")
	return &response.Member{
		ID:        member.ID,
		UUID:      member.UUID,
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
	}, nil
}

func (uc *member) DeleteMember(c *gin.Context, uuid string) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "DeleteMember called")
	if err := uc.memberRepo.DeleteMember(c, uuid); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to delete member: "+err.Error())
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "DeleteMember succeeded")
	return nil
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
