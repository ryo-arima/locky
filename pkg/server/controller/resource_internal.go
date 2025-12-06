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
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/server/share"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

// Local aliases for cleaner logging code
var (
	INFO  = share.GetServerLogger().INFO
	DEBUG = share.GetServerLogger().DEBUG
	WARN  = share.GetServerLogger().WARN
	ERROR = share.GetServerLogger().ERROR
)

// Local MCode definitions
var (
	SRNRSR1 = global.SRNRSR1
	SRNRSR2 = global.SRNRSR2
	Mcode   = global.Mcode
)

// ResourceInternal provides authenticated resource operations for internal scope.
type ResourceInternal interface {
	GetResources(c *gin.Context)
	CreateResource(c *gin.Context)
	UpdateResource(c *gin.Context)
	DeleteResource(c *gin.Context)
	CountResources(c *gin.Context)
}

type resourceInternal struct {
	ResourceUsecase usecase.Resource
}

func NewResourceInternal(resourceUsecase usecase.Resource) ResourceInternal {
	return &resourceInternal{
		ResourceUsecase: resourceUsecase,
	}
}

// swagger:operation GET /internal/resources resources getResourcesInternal
// ---
// summary: Get a list of resources the user has access to.
// description: Get a list of resources filtered by user's group memberships.
// responses:
//   "200":
//     description: A list of resources.
//   "400":
//     description: Bad request.
//   "401":
//     description: Unauthorized.
func (rcvr *resourceInternal) GetResources(c *gin.Context) {
	INFO(Mcode(SRNRSR1), "GetResources called")
	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		ERROR(Mcode(SRNRSR2), "Failed to bind request")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	// Get user UUID from JWT token
	userUUID, ok := share.GetUserUUID(c)
	if !ok {
		ERROR(Mcode(SRNRSR2), "User not authenticated")
		c.JSON(http.StatusUnauthorized, &response.Resources{Code: "UNAUTHORIZED", Message: "User not authenticated", Resources: []response.Resource{}})
		return
	}

	// Build filter; access filtering applied in usecase
	filter := model.ResourceQueryFilter{}

	// Apply additional query filters
	if id := c.Query("id"); id != "" {
		resource, err := rcvr.ResourceUsecase.GetResourceAccessible(c, userUUID, id)
		if err != nil {
			if err.Error() == "access denied" {
				ERROR(Mcode(SRNRSR2), "Access denied")
				c.JSON(http.StatusForbidden, &response.Resources{Code: "ACCESS_DENIED", Message: "User does not have access to this resource", Resources: []response.Resource{}})
				return
			}
			ERROR(Mcode(SRNRSR2), "Resource not found")
			c.JSON(http.StatusNotFound, &response.Resources{Code: "RESOURCE_NOT_FOUND", Message: err.Error(), Resources: []response.Resource{}})
			return
		}
		resp := response.Resource{
			ID:        resource.ID,
			UUID:      resource.UUID,
			Name:      resource.Name,
			Type:      resource.Type,
			GroupUUID: resource.GroupUUID,
			CreatedAt: resource.CreatedAt,
			UpdatedAt: resource.UpdatedAt,
			DeletedAt: resource.DeletedAt,
		}
		INFO(Mcode(SRNRSR1), "GetResources succeeded")
		c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource retrieved successfully", Resources: []response.Resource{resp}})
		return
	}

	if v := c.Query("type"); v != "" {
		filter.Type = &v
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
	if v := c.Query("group_uuid"); v != "" {
		filter.GroupUUID = &v
	}
	if v := c.Query("limit"); v != "" {
		if lim, err := strconv.Atoi(v); err == nil {
			filter.Limit = lim
		}
	}
	if v := c.Query("offset"); v != "" {
		if off, err := strconv.Atoi(v); err == nil {
			filter.Offset = off
		}
	}

	resources, err := rcvr.ResourceUsecase.GetResourcesAccessible(c, userUUID, filter)
	if err != nil {
		ERROR(Mcode(SRNRSR2), "Failed to list resources")
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_LIST_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resp := make([]response.Resource, 0, len(resources))
	for _, r := range resources {
		resp = append(resp, response.Resource{
			ID:        r.ID,
			UUID:      r.UUID,
			Name:      r.Name,
			Type:      r.Type,
			GroupUUID: r.GroupUUID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			DeletedAt: r.DeletedAt,
		})
	}
	INFO(Mcode(SRNRSR1), "GetResources succeeded")
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resources retrieved successfully", Resources: resp})
}

// swagger:operation POST /internal/resource resource createResourceInternal
// ---
// summary: Create a new resource.
// description: Create a new resource in a group the user has editor or owner access to.
// responses:
//   "200":
//     description: Resource created.
//   "400":
//     description: Bad request.
//   "401":
//     description: Unauthorized.
//   "403":
//     description: Forbidden - user does not have editor/owner access.
func (rcvr *resourceInternal) CreateResource(c *gin.Context) {
	INFO(Mcode(SRNRSR1), "CreateResource called")
	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		ERROR(Mcode(SRNRSR2), "Failed to bind request")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	if resourceRequest.Name == "" {
		ERROR(Mcode(SRNRSR2), "Resource name is required")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_NAME_REQUIRED", Message: "Resource name is required", Resources: []response.Resource{}})
		return
	}

	if resourceRequest.GroupUUID == "" {
		ERROR(Mcode(SRNRSR2), "Group UUID is required")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "GROUP_UUID_REQUIRED", Message: "Group UUID is required", Resources: []response.Resource{}})
		return
	}

	// Get user UUID from JWT token
	userUUID, ok := share.GetUserUUID(c)
	if !ok {
		ERROR(Mcode(SRNRSR2), "User not authenticated")
		c.JSON(http.StatusUnauthorized, &response.Resources{Code: "UNAUTHORIZED", Message: "User not authenticated", Resources: []response.Resource{}})
		return
	}

	now := time.Now()
	resource := &model.Resources{
		UUID:      uuid.New().String(),
		Name:      resourceRequest.Name,
		Type:      resourceRequest.Type,
		GroupUUID: resourceRequest.GroupUUID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := rcvr.ResourceUsecase.CreateResourceWithAccess(c, userUUID, resource); err != nil {
		if err.Error() == "access denied" || err.Error() == "insufficient role to create resource" {
			ERROR(Mcode(SRNRSR2), "Access denied")
			c.JSON(http.StatusForbidden, &response.Resources{Code: "ACCESS_DENIED", Message: err.Error(), Resources: []response.Resource{}})
			return
		}
		ERROR(Mcode(SRNRSR2), "Failed to create resource")
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_CREATE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resp := response.Resource{
		ID:        resource.ID,
		UUID:      resource.UUID,
		Name:      resource.Name,
		Type:      resource.Type,
		GroupUUID: resource.GroupUUID,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
	}
	INFO(Mcode(SRNRSR1), "CreateResource succeeded")
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource created successfully", Resources: []response.Resource{resp}})
}

// swagger:operation PUT /internal/resource/:id resource updateResourceInternal
// ---
// summary: Update an existing resource.
// description: Update resource details. Requires editor or owner access to the resource's group.
// responses:
//   "200":
//     description: Resource updated.
//   "400":
//     description: Bad request.
//   "401":
//     description: Unauthorized.
//   "403":
//     description: Forbidden - user does not have editor/owner access.
func (rcvr *resourceInternal) UpdateResource(c *gin.Context) {
	INFO(Mcode(SRNRSR1), "UpdateResource called")
	id := c.Param("id")
	if id == "" {
		ERROR(Mcode(SRNRSR2), "Resource ID is required")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_ID_REQUIRED", Message: "Resource ID is required", Resources: []response.Resource{}})
		return
	}

	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		ERROR(Mcode(SRNRSR2), "Failed to bind request")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	// Get user UUID from JWT token
	userUUID, ok := share.GetUserUUID(c)
	if !ok {
		ERROR(Mcode(SRNRSR2), "User not authenticated")
		c.JSON(http.StatusUnauthorized, &response.Resources{Code: "UNAUTHORIZED", Message: "User not authenticated", Resources: []response.Resource{}})
		return
	}

	// Retrieve resource and ensure user has access (usecase will validate)
	resource, err := rcvr.ResourceUsecase.GetResourceAccessible(c, userUUID, id)
	if err != nil {
		if err.Error() == "access denied" {
			ERROR(Mcode(SRNRSR2), "Access denied")
			c.JSON(http.StatusForbidden, &response.Resources{Code: "ACCESS_DENIED", Message: "User does not have access to this resource", Resources: []response.Resource{}})
			return
		}
		ERROR(Mcode(SRNRSR2), "Resource not found")
		c.JSON(http.StatusNotFound, &response.Resources{Code: "RESOURCE_NOT_FOUND", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	now := time.Now()
	if resourceRequest.Name != "" {
		resource.Name = resourceRequest.Name
	}
	if resourceRequest.Type != "" {
		resource.Type = resourceRequest.Type
	}
	resource.UpdatedAt = &now

	if err := rcvr.ResourceUsecase.UpdateResourceWithAccess(c, userUUID, resource); err != nil {
		if err.Error() == "access denied" || err.Error() == "insufficient role to update resource" {
			ERROR(Mcode(SRNRSR2), "Access denied")
			c.JSON(http.StatusForbidden, &response.Resources{Code: "ACCESS_DENIED", Message: err.Error(), Resources: []response.Resource{}})
			return
		}
		ERROR(Mcode(SRNRSR2), "Failed to update resource")
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_UPDATE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resp := response.Resource{
		ID:        resource.ID,
		UUID:      resource.UUID,
		Name:      resource.Name,
		Type:      resource.Type,
		GroupUUID: resource.GroupUUID,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
	}
	INFO(Mcode(SRNRSR1), "UpdateResource succeeded")
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource updated successfully", Resources: []response.Resource{resp}})
}

// swagger:operation DELETE /internal/resource/:id resource deleteResourceInternal
// ---
// summary: Delete a resource.
// description: Soft delete a resource by UUID.
// responses:
//   "200":
//     description: Resource deleted.
//   "400":
//     description: Bad request.
func (rcvr *resourceInternal) DeleteResource(c *gin.Context) {
	INFO(Mcode(SRNRSR1), "DeleteResource called")
	id := c.Param("id")
	if id == "" {
		ERROR(Mcode(SRNRSR2), "Resource ID is required")
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_ID_REQUIRED", Message: "Resource ID is required", Resources: []response.Resource{}})
		return
	}

	// Get user UUID and delegate access check to usecase
	userUUID, ok := share.GetUserUUID(c)
	if !ok {
		ERROR(Mcode(SRNRSR2), "User not authenticated")
		c.JSON(http.StatusUnauthorized, &response.Resources{Code: "UNAUTHORIZED", Message: "User not authenticated", Resources: []response.Resource{}})
		return
	}

	if err := rcvr.ResourceUsecase.DeleteResourceWithAccess(c, userUUID, id); err != nil {
		if err.Error() == "access denied" || err.Error() == "insufficient role to delete resource" {
			ERROR(Mcode(SRNRSR2), "Access denied")
			c.JSON(http.StatusForbidden, &response.Resources{Code: "ACCESS_DENIED", Message: err.Error(), Resources: []response.Resource{}})
			return
		}
		ERROR(Mcode(SRNRSR2), "Failed to delete resource")
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_DELETE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	INFO(Mcode(SRNRSR1), "DeleteResource succeeded")
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource deleted successfully", Resources: []response.Resource{}})
}

// swagger:operation GET /internal/resources/count resource countResourcesInternal
// ---
// summary: Count resources.
// description: Get the total count of resources.
// responses:
//   "200":
//     description: Resource count.
func (rcvr *resourceInternal) CountResources(c *gin.Context) {
	INFO(Mcode(SRNRSR1), "CountResources called")
	// Get user UUID
	userUUID, ok := share.GetUserUUID(c)
	if !ok {
		ERROR(Mcode(SRNRSR2), "User not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "User not authenticated"})
		return
	}

	// Build filter from query params
	var filter model.ResourceQueryFilter
	if v := c.Query("type"); v != "" {
		filter.Type = &v
	}
	if v := c.Query("group_uuid"); v != "" {
		filter.GroupUUID = &v
	}

	count, err := rcvr.ResourceUsecase.CountResourcesAccessible(c, userUUID, filter)
	if err != nil {
		ERROR(Mcode(SRNRSR2), "Failed to count resources")
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_COUNT_ERROR", "message": err.Error()})
		return
	}
	INFO(Mcode(SRNRSR1), "CountResources succeeded")
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "Resource count retrieved", "count": count})
}
