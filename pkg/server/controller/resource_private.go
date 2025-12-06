package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/usecase"
)

// ResourcePrivate provides administrative resource operations for private scope.
type ResourcePrivate interface {
	GetResources(c *gin.Context)
	CreateResource(c *gin.Context)
	UpdateResource(c *gin.Context)
	DeleteResource(c *gin.Context)
	CountResources(c *gin.Context)
}

type resourcePrivate struct {
	ResourceUsecase usecase.Resource
}

func NewResourcePrivate(resourceUsecase usecase.Resource) ResourcePrivate {
	return &resourcePrivate{
		ResourceUsecase: resourceUsecase,
	}
}

// swagger:operation GET /private/resources resources getResourcesPrivate
// ---
// summary: Get a list of resources (admin).
// description: Get a list of all resources in the system (administrative access).
// responses:
//   "200":
//     description: A list of resources.
//   "400":
//     description: Bad request.
func (rcvr *resourcePrivate) GetResources(c *gin.Context) {
	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	if id := c.Query("id"); id != "" {
		resource, err := rcvr.ResourceUsecase.GetResourceAccessibleAdmin(c, "", id)
		if err != nil {
			c.JSON(http.StatusNotFound, &response.Resources{Code: "RESOURCE_NOT_FOUND", Message: err.Error(), Resources: []response.Resource{}})
			return
		}
		resp := response.Resource{
			ID:        resource.ID,
			UUID:      resource.UUID,
			Name:      resource.Name,
			Type:      resource.Type,
			CreatedAt: resource.CreatedAt,
			UpdatedAt: resource.UpdatedAt,
			DeletedAt: resource.DeletedAt,
		}
		c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource retrieved successfully", Resources: []response.Resource{resp}})
		return
	}

	resources, err := rcvr.ResourceUsecase.GetResourcesAccessibleAdmin(c, "", repository.ResourceQueryFilter{})
	if err != nil {
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
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			DeletedAt: r.DeletedAt,
		})
	}
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resources retrieved successfully", Resources: resp})
}

// swagger:operation POST /private/resource resource createResourcePrivate
// ---
// summary: Create a new resource (admin).
// description: Create a new resource in the system with administrative privileges.
// responses:
//   "200":
//     description: Resource created.
//   "400":
//     description: Bad request.
func (rcvr *resourcePrivate) CreateResource(c *gin.Context) {
	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	if resourceRequest.Name == "" {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_NAME_REQUIRED", Message: "Resource name is required", Resources: []response.Resource{}})
		return
	}

	now := time.Now()
	resource := &model.Resources{
		UUID:      uuid.New().String(),
		Name:      resourceRequest.Name,
		Type:      resourceRequest.Type,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := rcvr.ResourceUsecase.CreateResourceWithAdmin(c, "", resource); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_CREATE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resp := response.Resource{
		ID:        resource.ID,
		UUID:      resource.UUID,
		Name:      resource.Name,
		Type:      resource.Type,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
	}
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource created successfully", Resources: []response.Resource{resp}})
}

// swagger:operation PUT /private/resource/:id resource updateResourcePrivate
// ---
// summary: Update an existing resource (admin).
// description: Update resource details by UUID with administrative privileges.
// responses:
//   "200":
//     description: Resource updated.
//   "400":
//     description: Bad request.
func (rcvr *resourcePrivate) UpdateResource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_ID_REQUIRED", Message: "Resource ID is required", Resources: []response.Resource{}})
		return
	}

	var resourceRequest request.Resource
	if err := c.Bind(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_BIND_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resource, err := rcvr.ResourceUsecase.GetResourceAccessibleAdmin(c, "", id)
	if err != nil {
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

	if err := rcvr.ResourceUsecase.UpdateResourceWithAdmin(c, "", resource); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_UPDATE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	resp := response.Resource{
		ID:        resource.ID,
		UUID:      resource.UUID,
		Name:      resource.Name,
		Type:      resource.Type,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
	}
	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource updated successfully", Resources: []response.Resource{resp}})
}

// swagger:operation DELETE /private/resource/:id resource deleteResourcePrivate
// ---
// summary: Delete a resource (admin).
// description: Soft delete a resource by UUID with administrative privileges.
// responses:
//   "200":
//     description: Resource deleted.
//   "400":
//     description: Bad request.
func (rcvr *resourcePrivate) DeleteResource(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &response.Resources{Code: "RESOURCE_ID_REQUIRED", Message: "Resource ID is required", Resources: []response.Resource{}})
		return
	}

	if err := rcvr.ResourceUsecase.DeleteResourceWithAdmin(c, "", id); err != nil {
		c.JSON(http.StatusInternalServerError, &response.Resources{Code: "RESOURCE_DELETE_ERROR", Message: err.Error(), Resources: []response.Resource{}})
		return
	}

	c.JSON(http.StatusOK, &response.Resources{Code: "SUCCESS", Message: "Resource deleted successfully", Resources: []response.Resource{}})
}

// swagger:operation GET /private/resources/count resource countResourcesPrivate
// ---
// summary: Count resources (admin).
// description: Get the total count of resources with administrative access.
// responses:
//   "200":
//     description: Resource count.
func (rcvr *resourcePrivate) CountResources(c *gin.Context) {
	count, err := rcvr.ResourceUsecase.CountResourcesAccessibleAdmin(c, "", repository.ResourceQueryFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "RESOURCE_COUNT_ERROR", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "SUCCESS", "message": "Resource count retrieved", "count": count})
}
