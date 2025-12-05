package usecase

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"gorm.io/gorm"
)

type Group interface {
	GetGroups(c *gin.Context) ([]response.Group, error)
	GetGroupByUUID(c *gin.Context, uuid string) (*response.Group, error)
	GetGroupByID(c *gin.Context, id uint) (*response.Group, error)
	CreateGroup(c *gin.Context, group *model.Groups) (*response.Group, *gorm.DB)
	UpdateGroup(c *gin.Context, group *model.Groups) (*response.Group, *gorm.DB)
	DeleteGroup(c *gin.Context, uuid string) *gorm.DB
	ListGroups(c *gin.Context, filter repository.GroupQueryFilter) ([]response.Group, error)
	CountGroups(c *gin.Context, filter repository.GroupQueryFilter) (int64, error)
}

type group struct {
	groupRepo repository.Group
}

func NewGroup(groupRepo repository.Group) Group {
	return &group{
		groupRepo: groupRepo,
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

func (uc *group) CreateGroup(c *gin.Context, group *model.Groups) (*response.Group, *gorm.DB) {
	resDB := uc.groupRepo.CreateGroup(c, group)
	if resDB.Error != nil {
		return nil, resDB
	}
	resp := convertGroupModelToResponse(*group)
	return &resp, resDB
}

func (uc *group) UpdateGroup(c *gin.Context, group *model.Groups) (*response.Group, *gorm.DB) {
	resDB := uc.groupRepo.UpdateGroup(c, group)
	if resDB.Error != nil {
		return nil, resDB
	}
	resp := convertGroupModelToResponse(*group)
	return &resp, resDB
}

func (uc *group) DeleteGroup(c *gin.Context, uuid string) *gorm.DB {
	return uc.groupRepo.DeleteGroup(c, uuid)
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
