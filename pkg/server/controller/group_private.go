package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/middleware"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type GroupControllerForPrivate interface {
	GetGroups(c *gin.Context)
	CreateGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
	CountGroups(c *gin.Context)
}

type groupControllerForPrivate struct {
	GroupRepository  repository.GroupRepository
	CommonRepository repository.CommonRepository
}

func (groupController groupControllerForPrivate) GetGroups(c *gin.Context) {
	// swagger:operation GET /private/groups groups getGroupsPrivate
	// ---
	// summary: Get a list of groups.
	// description: Get a list of all groups in the system.
	// responses:
	//   "200":
	//     description: A list of groups.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	filter := repository.GroupQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("name"); v != "" {
		filter.Name = &v
	}
	if v := c.Query("name_prefix"); v != "" {
		filter.NamePrefix = &v
	}
	if v := c.Query("name_like"); v != "" {
		filter.NameLike = &v
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Limit = n
		}
	}
	if v := c.Query("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.Offset = n
		}
	}
	groups, err := groupController.GroupRepository.ListGroups(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.GroupResponse{Code: "SERVER_CONTROLLER_GET__FOR__002", Message: err.Error(), Groups: []response.Group{}})
		return
	}
	resp := make([]response.Group, 0, len(groups))
	for _, g := range groups {
		resp = append(resp, response.Group{ID: g.ID, UUID: g.UUID, Name: g.Name, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt, DeletedAt: g.DeletedAt})
	}
	c.JSON(http.StatusOK, &response.GroupResponse{Code: "SUCCESS", Message: "Groups retrieved successfully", Groups: resp})
}

func (groupController groupControllerForPrivate) CountGroups(c *gin.Context) {
	filter := repository.GroupQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("name"); v != "" {
		filter.Name = &v
	}
	if v := c.Query("name_prefix"); v != "" {
		filter.NamePrefix = &v
	}
	if v := c.Query("name_like"); v != "" {
		filter.NameLike = &v
	}
	cnt, err := groupController.GroupRepository.CountGroups(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "SERVER_CONTROLLER_COUNT__FOR__001", "message": err.Error(), "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "Count retrieved", "count": cnt})
}

func (groupController groupControllerForPrivate) CreateGroup(c *gin.Context) {
	// swagger:operation POST /private/groups groups createGroupPrivate
	// ---
	// summary: Create a new group.
	// description: Create a new group with the provided information.
	// parameters:
	// - name: group
	//   in: body
	//   description: The group to create.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/GroupRequest"
	// responses:
	//   "200":
	//     description: The created group.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	var groupRequest request.GroupRequest
	if err := c.Bind(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__001", Message: err.Error(), Groups: []response.Group{}})
		return
	}
	if groupRequest.Name == "" {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: "name is required", Groups: []response.Group{}})
		return
	}
	now := time.Now()
	g := model.Groups{UUID: uuid.New().String(), Name: groupRequest.Name, CreatedAt: &now, UpdatedAt: &now}
	resDB := groupController.GroupRepository.CreateGroup(&g)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.GroupResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__003", Message: resDB.Error.Error(), Groups: []response.Group{}})
		return
	}
	claims, ok := middleware.GetUserClaims(c)
	if ok && claims != nil {
		memberRepo := repository.NewMemberRepository(groupController.CommonRepository.GetBaseConfig())
		mem := model.Members{UUID: uuid.New().String(), GroupUUID: g.UUID, UserUUID: claims.UUID, Role: "owner", CreatedAt: &now, UpdatedAt: &now}
		_ = memberRepo.CreateMember(&mem)
	}
	c.JSON(http.StatusOK, &response.GroupResponse{Code: "SUCCESS", Message: "Group created successfully", Groups: []response.Group{{ID: g.ID, UUID: g.UUID, Name: g.Name}}})
}

func (groupController groupControllerForPrivate) UpdateGroup(c *gin.Context) {
	// swagger:operation PUT /private/groups/{id} groups updateGroupPrivate
	// ---
	// summary: Update a group.
	// description: Update a group with the provided information.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the group to update.
	//   required: true
	//   type: integer
	// - name: group
	//   in: body
	//   description: The group to update.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/GroupRequest"
	// responses:
	//   "200":
	//     description: The updated group.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	var groupRequest request.GroupRequest
	if err := c.Bind(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__001", Message: err.Error(), Groups: []response.Group{}})
		return
	}
	if groupRequest.ID == 0 {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__002", Message: "id is required", Groups: []response.Group{}})
		return
	}
	g, err := groupController.GroupRepository.GetGroupByID(groupRequest.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__003", Message: "group not found", Groups: []response.Group{}})
		return
	}
	if groupRequest.Name != "" {
		g.Name = groupRequest.Name
	}
	now := time.Now()
	g.UpdatedAt = &now
	resDB := groupController.GroupRepository.UpdateGroup(&g)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.GroupResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__004", Message: resDB.Error.Error(), Groups: []response.Group{}})
		return
	}
	c.JSON(http.StatusOK, &response.GroupResponse{Code: "SUCCESS", Message: "Group updated successfully", Groups: []response.Group{{ID: g.ID, UUID: g.UUID, Name: g.Name}}})
}

func (groupController groupControllerForPrivate) DeleteGroup(c *gin.Context) {
	// swagger:operation DELETE /private/groups/{id} groups deleteGroupPrivate
	// ---
	// summary: Delete a group.
	// description: Delete a group by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the group to delete.
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: The deleted group.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/GroupResponse"
	var groupRequest request.GroupRequest
	if err := c.Bind(&groupRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__001", Message: err.Error(), Groups: []response.Group{}})
		return
	}
	if groupRequest.UUID == "" {
		c.JSON(http.StatusBadRequest, &response.GroupResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__002", Message: "uuid is required", Groups: []response.Group{}})
		return
	}
	resDB := groupController.GroupRepository.DeleteGroup(groupRequest.UUID)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.GroupResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__003", Message: resDB.Error.Error(), Groups: []response.Group{}})
		return
	}
	c.JSON(http.StatusOK, &response.GroupResponse{Code: "SUCCESS", Message: "Group deleted successfully", Groups: []response.Group{}})
}

func NewGroupControllerForPrivate(groupRepository repository.GroupRepository, commonRepository repository.CommonRepository) GroupControllerForPrivate {
	return &groupControllerForPrivate{GroupRepository: groupRepository, CommonRepository: commonRepository}
}
