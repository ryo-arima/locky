package usecase

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
	"gorm.io/gorm"
)

type Group interface {
	GetGroups(c *gin.Context) ([]response.Group, error)
	GetGroupByUUID(c *gin.Context, uuid string) (*response.Group, error)
	GetGroupByID(c *gin.Context, id uint) (*response.Group, error)
	CreateGroup(c *gin.Context, group *model.Groups) (*response.Group, error)
	UpdateGroup(c *gin.Context, group *model.Groups) (*response.Group, error)
	DeleteGroup(c *gin.Context, uuid string) error
	ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]response.Group, error)
	CountGroups(c *gin.Context, filter repository.GroupQueryFilter) (int64, error)
	// CreateGroupWithOwnerMember atomically creates a group and registers the
	// specified user as the owner member within a single DB transaction.
	CreateGroupWithOwnerMember(c *gin.Context, group *model.Groups, ownerUserUUID string) (*response.Group, error)
}

type group struct {
	groupRepo  repository.Group
	memberRepo repository.Member
	db         *gorm.DB
}

func NewGroup(groupRepo repository.Group, memberRepo repository.Member, db *gorm.DB) Group {
	return &group{
		groupRepo:  groupRepo,
		memberRepo: memberRepo,
		db:         db,
	}
}

// convertGroupModelToResponse converts a model.Groups to response.Group
func convertGroupModelToResponse(group model.Groups) response.Group {
	return response.Group{
		ID:   group.ID,
		UUID: group.UUID,
		Name: group.Name,
	}
}

// convertGroupModelsToResponses converts []model.Groups to []response.Group
func convertGroupModelsToResponses(groups []model.Groups) []response.Group {
	responseGroups := make([]response.Group, 0, len(groups))
	for _, group := range groups {
		responseGroups = append(responseGroups, convertGroupModelToResponse(group))
	}
	return responseGroups
}

func (uc *group) GetGroups(c *gin.Context) ([]response.Group, error) {
	groups := uc.groupRepo.GetGroups(c)
	return convertGroupModelsToResponses(groups), nil
}

func (uc *group) GetGroupByUUID(c *gin.Context, uuid string) (*response.Group, error) {
	group, err := uc.groupRepo.GetGroupByUUID(c, uuid)
	if err != nil {
		return nil, err
	}
	resp := convertGroupModelToResponse(group)
	return &resp, nil
}

func (uc *group) GetGroupByID(c *gin.Context, id uint) (*response.Group, error) {
	group, err := uc.groupRepo.GetGroupByID(c, id)
	if err != nil {
		return nil, err
	}
	resp := convertGroupModelToResponse(group)
	return &resp, nil
}

func (uc *group) CreateGroup(c *gin.Context, group *model.Groups) (*response.Group, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CreateGroup called")
	if err := uc.groupRepo.CreateGroup(c, group); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to create group: "+err.Error())
		return nil, err
	}
	resp := convertGroupModelToResponse(*group)
	INFO(reqID, Mcode(SRNRSR1), "CreateGroup succeeded")
	return &resp, nil
}

func (uc *group) UpdateGroup(c *gin.Context, group *model.Groups) (*response.Group, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "UpdateGroup called")
	if err := uc.groupRepo.UpdateGroup(c, group); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to update group: "+err.Error())
		return nil, err
	}
	resp := convertGroupModelToResponse(*group)
	INFO(reqID, Mcode(SRNRSR1), "UpdateGroup succeeded")
	return &resp, nil
}

func (uc *group) DeleteGroup(c *gin.Context, uuid string) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "DeleteGroup called")
	if err := uc.groupRepo.DeleteGroup(c, uuid); err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to delete group: "+err.Error())
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "DeleteGroup succeeded")
	return nil
}

// CreateGroupWithOwnerMember atomically creates a group and an owner member in one transaction.
// If either operation fails, the entire transaction is rolled back.
func (uc *group) CreateGroupWithOwnerMember(c *gin.Context, g *model.Groups, ownerUserUUID string) (*response.Group, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CreateGroupWithOwnerMember called")

	var resp response.Group
	err := repository.RunInTx(uc.db, func(tx *gorm.DB) error {
		// Create group within transaction
		if err := uc.groupRepo.WithTx(tx).CreateGroup(c, g); err != nil {
			ERROR(reqID, Mcode(SRNRSR2), "Failed to create group in tx: "+err.Error())
			return err
		}

		// Register owner member within the same transaction
		now := time.Now()
		mem := &model.Members{
			UUID:      uuid.New().String(),
			GroupUUID: g.UUID,
			UserUUID:  ownerUserUUID,
			Role:      "owner",
			CreatedAt: &now,
			UpdatedAt: &now,
		}
		if err := uc.memberRepo.WithTx(tx).CreateMember(c, mem); err != nil {
			ERROR(reqID, Mcode(SRNRSR2), "Failed to create owner member in tx: "+err.Error())
			return err
		}

		resp = convertGroupModelToResponse(*g)
		return nil
	})
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "CreateGroupWithOwnerMember transaction failed: "+err.Error())
		return nil, err
	}

	INFO(reqID, Mcode(SRNRSR1), "CreateGroupWithOwnerMember succeeded")
	return &resp, nil
}

func (uc *group) ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]response.Group, error) {
	groups, err := uc.groupRepo.ListGroups(c, filter)
	if err != nil {
		return nil, err
	}
	return convertGroupModelsToResponses(groups), nil
}

func (uc *group) CountGroups(c *gin.Context, filter repository.GroupQueryFilter) (int64, error) {
	return uc.groupRepo.CountGroups(c, filter)
}
