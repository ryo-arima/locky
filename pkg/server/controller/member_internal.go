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
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// MemberControllerForInternal provides authenticated member operations for internal scope.
//
// Exposes CRUD endpoints requiring a valid bearer token:
//   - GetMembers: List members (GET /v1/internal/members)
//   - CreateMember: Create member (POST /v1/internal/members)
//   - UpdateMember: Update member (PUT /v1/internal/members/{id})
//   - DeleteMember: Delete member (DELETE /v1/internal/members/{id})
type MemberControllerForInternal interface {
	GetMembers(c *gin.Context)
	CreateMember(c *gin.Context)
	UpdateMember(c *gin.Context)
	DeleteMember(c *gin.Context)
	CountMembers(c *gin.Context)
}

type memberControllerForInternal struct {
	MemberRepository repository.MemberRepository
	CommonRepository repository.CommonRepository
}

// GetMembers lists members (authenticated).
//
// Route: GET /v1/internal/members
// Security: Bearer token
func (memberController memberControllerForInternal) GetMembers(c *gin.Context) {
	filter := repository.MemberQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("group_uuid"); v != "" {
		filter.GroupUUID = &v
	}
	if v := c.Query("user_uuid"); v != "" {
		filter.UserUUID = &v
	}
	if v := c.Query("role"); v != "" {
		filter.Role = &v
	}
	if v := c.Query("role_prefix"); v != "" {
		filter.RolePrefix = &v
	}
	if v := c.Query("role_like"); v != "" {
		filter.RoleLike = &v
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
	members, err := memberController.MemberRepository.ListMembers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &response.MemberResponse{Code: "SERVER_CONTROLLER_GET__FOR__002", Message: err.Error(), Members: []response.Member{}})
		return
	}
	resp := make([]response.Member, 0, len(members))
	for _, m := range members {
		resp = append(resp, response.Member{ID: m.ID, UUID: m.UUID, GroupUUID: m.GroupUUID, UserUUID: m.UserUUID, Role: m.Role, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt, DeletedAt: m.DeletedAt})
	}
	c.JSON(http.StatusOK, &response.MemberResponse{Code: "SUCCESS", Message: "Members retrieved successfully", Members: resp})
}

// CreateMember creates a member (authenticated).
//
// Route: POST /v1/internal/members
// Security: Bearer token
func (memberController memberControllerForInternal) CreateMember(c *gin.Context) {
	// swagger:operation POST /internal/members members createMemberInternal
	// ---
	// summary: Create a new member.
	// description: Create a new member with the provided information.
	// parameters:
	// - name: member
	//   in: body
	//   description: The member to create.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/MemberRequest"
	// responses:
	//   "200":
	//     description: The created member.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	var memberRequest request.MemberRequest
	if err := c.Bind(&memberRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__001", Message: err.Error(), Members: []response.Member{}})
		return
	}
	if memberRequest.GroupUUID == "" || memberRequest.UserUUID == "" || memberRequest.Role == "" {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__002", Message: "group_uuid, user_uuid and role are required", Members: []response.Member{}})
		return
	}
	now := time.Now()
	m := model.Members{UUID: uuid.New().String(), GroupUUID: memberRequest.GroupUUID, UserUUID: memberRequest.UserUUID, Role: memberRequest.Role, CreatedAt: &now, UpdatedAt: &now}
	resDB := memberController.MemberRepository.CreateMember(&m)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.MemberResponse{Code: "SERVER_CONTROLLER_CREATE__FOR__003", Message: resDB.Error.Error(), Members: []response.Member{}})
		return
	}
	c.JSON(http.StatusOK, &response.MemberResponse{Code: "SUCCESS", Message: "Member created successfully", Members: []response.Member{{ID: m.ID, UUID: m.UUID, GroupUUID: m.GroupUUID, UserUUID: m.UserUUID, Role: m.Role}}})
}

// UpdateMember updates a member (authenticated).
//
// Route: PUT /v1/internal/members/{id}
// Security: Bearer token
func (memberController memberControllerForInternal) UpdateMember(c *gin.Context) {
	// swagger:operation PUT /internal/members/{id} members updateMemberInternal
	// ---
	// summary: Update a member.
	// description: Update a member with the provided information.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the member to update.
	//   required: true
	//   type: integer
	// - name: member
	//   in: body
	//   description: The member to update.
	//   required: true
	//   schema:
	//     $ref: "#/definitions/MemberRequest"
	// responses:
	//   "200":
	//     description: The updated member.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	var memberRequest request.MemberRequest
	if err := c.Bind(&memberRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__001", Message: err.Error(), Members: []response.Member{}})
		return
	}
	if memberRequest.ID == 0 {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__002", Message: "id is required", Members: []response.Member{}})
		return
	}
	now := time.Now()
	upd := model.Members{ID: memberRequest.ID, GroupUUID: memberRequest.GroupUUID, UserUUID: memberRequest.UserUUID, Role: memberRequest.Role, UpdatedAt: &now}
	resDB := memberController.MemberRepository.UpdateMember(&upd)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.MemberResponse{Code: "SERVER_CONTROLLER_UPDATE__FOR__003", Message: resDB.Error.Error(), Members: []response.Member{}})
		return
	}
	c.JSON(http.StatusOK, &response.MemberResponse{Code: "SUCCESS", Message: "Member updated successfully", Members: []response.Member{{ID: upd.ID, UUID: upd.UUID, GroupUUID: upd.GroupUUID, UserUUID: upd.UserUUID, Role: upd.Role}}})
}

// DeleteMember deletes a member (authenticated).
//
// Route: DELETE /v1/internal/members/{id}
// Security: Bearer token
func (memberController memberControllerForInternal) DeleteMember(c *gin.Context) {
	// swagger:operation DELETE /internal/members/{id} members deleteMemberInternal
	// ---
	// summary: Delete a member.
	// description: Delete a member by ID.
	// parameters:
	// - name: id
	//   in: path
	//   description: The ID of the member to delete.
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: The deleted member.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	//   "400":
	//     description: Bad request.
	//     schema:
	//       $ref: "#/definitions/MemberResponse"
	var memberRequest request.MemberRequest
	if err := c.Bind(&memberRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__001", Message: err.Error(), Members: []response.Member{}})
		return
	}
	if memberRequest.UUID == "" {
		c.JSON(http.StatusBadRequest, &response.MemberResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__002", Message: "uuid is required", Members: []response.Member{}})
		return
	}
	resDB := memberController.MemberRepository.DeleteMember(memberRequest.UUID)
	if resDB.Error != nil {
		c.JSON(http.StatusInternalServerError, &response.MemberResponse{Code: "SERVER_CONTROLLER_DELETE__FOR__003", Message: resDB.Error.Error(), Members: []response.Member{}})
		return
	}
	c.JSON(http.StatusOK, &response.MemberResponse{Code: "SUCCESS", Message: "Member deleted successfully", Members: []response.Member{}})
}

// CountMembers counts members (authenticated).
//
// Route: GET /v1/internal/members/count
// Security: Bearer token
func (memberController memberControllerForInternal) CountMembers(c *gin.Context) {
	filter := repository.MemberQueryFilter{}
	if v := c.Query("id"); v != "" {
		if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
			id := uint(id64)
			filter.ID = &id
		}
	}
	if v := c.Query("uuid"); v != "" {
		filter.UUID = &v
	}
	if v := c.Query("group_uuid"); v != "" {
		filter.GroupUUID = &v
	}
	if v := c.Query("user_uuid"); v != "" {
		filter.UserUUID = &v
	}
	if v := c.Query("role"); v != "" {
		filter.Role = &v
	}
	if v := c.Query("role_prefix"); v != "" {
		filter.RolePrefix = &v
	}
	if v := c.Query("role_like"); v != "" {
		filter.RoleLike = &v
	}
	cnt, err := memberController.MemberRepository.CountMembers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "SERVER_CONTROLLER_COUNT__FOR__001", "message": err.Error(), "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "Count retrieved", "count": cnt})
}

// NewMemberControllerForInternal creates a new internal member controller.
//
// Parameters:
//   - memberRepository: Member data repository
//   - commonRepository: Common services repository
//
// Returns:
//   - MemberControllerForInternal: Configured internal controller instance
func NewMemberControllerForInternal(memberRepository repository.MemberRepository, commonRepository repository.CommonRepository) MemberControllerForInternal {
	return &memberControllerForInternal{MemberRepository: memberRepository, CommonRepository: commonRepository}
}
