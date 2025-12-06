package usecase

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
	"github.com/ryo-arima/locky/pkg/server/share"
)

// Resource defines usecase operations for resource management with access checks.
type Resource interface {
	// List resources the user has access to according to their memberships.
	GetResourcesAccessible(c *gin.Context, userUUID string, filter model.ResourceQueryFilter) ([]model.Resources, error)
	// Get single resource and verify the user has access to it.
	GetResourceAccessible(c *gin.Context, userUUID string, id string) (*model.Resources, error)
	// Count resources the user has access to.
	CountResourcesAccessible(c *gin.Context, userUUID string, filter model.ResourceQueryFilter) (int64, error)
	// Create resource if user has editor/owner role on the target group.
	CreateResourceWithAccess(c *gin.Context, userUUID string, resource *model.Resources) error
	// Update resource if user has editor/owner role on the owning group.
	UpdateResourceWithAccess(c *gin.Context, userUUID string, resource *model.Resources) error
	// Delete resource if user has editor/owner role on the owning group.
	DeleteResourceWithAccess(c *gin.Context, userUUID string, id string) error
	// Helper: get group UUIDs the user is a member of (any role).
	GetAccessibleGroupUUIDs(c *gin.Context, userUUID string) ([]string, error)
	// Administrative passthroughs (no access checks) for private/admin controllers
	GetResourcesAccessibleAdmin(c *gin.Context, _ string, filter model.ResourceQueryFilter) ([]model.Resources, error)
	GetResourceAccessibleAdmin(c *gin.Context, _ string, id string) (*model.Resources, error)
	CountResourcesAccessibleAdmin(c *gin.Context, _ string, filter model.ResourceQueryFilter) (int64, error)
	CreateResourceWithAdmin(c *gin.Context, _ string, resource *model.Resources) error
	UpdateResourceWithAdmin(c *gin.Context, _ string, resource *model.Resources) error
	DeleteResourceWithAdmin(c *gin.Context, _ string, id string) error
}

type resource struct {
	resourceRepo repository.Resource
	memberRepo   repository.Member
}

func NewResource(resourceRepo repository.Resource, memberRepo repository.Member) Resource {
	return &resource{
		resourceRepo: resourceRepo,
		memberRepo:   memberRepo,
	}
}

// GetAccessibleGroupUUIDs returns group UUIDs where the user has a membership (any role).
func (uc *resource) GetAccessibleGroupUUIDs(c *gin.Context, userUUID string) ([]string, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "GetAccessibleGroupUUIDs called")
	filter := repository.MemberQueryFilter{
		UserUUID: &userUUID,
		Limit:    200,
	}
	members, err := uc.memberRepo.ListMembers(c, filter)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to list members")
		return nil, err
	}

	groupUUIDs := make([]string, 0, len(members))
	for _, m := range members {
		groupUUIDs = append(groupUUIDs, m.GroupUUID)
	}
	INFO(reqID, Mcode(SRNRSR1), "GetAccessibleGroupUUIDs succeeded")
	return groupUUIDs, nil
}

// helper: get member entries for specific user+group
func (uc *resource) getMemberForUserAndGroup(c *gin.Context, userUUID, groupUUID string) (*model.Members, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "getMemberForUserAndGroup called")
	filter := repository.MemberQueryFilter{
		UserUUID:  &userUUID,
		GroupUUID: &groupUUID,
		Limit:     1,
	}
	members, err := uc.memberRepo.ListMembers(c, filter)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to list members")
		return nil, err
	}
	if len(members) == 0 {
		INFO(reqID, Mcode(SRNRSR1), "No member found")
		return nil, nil
	}
	INFO(reqID, Mcode(SRNRSR1), "getMemberForUserAndGroup succeeded")
	return &members[0], nil
}

// GetResourcesAccessible returns resources filtered to only groups the user is a member of.
func (uc *resource) GetResourcesAccessible(c *gin.Context, userUUID string, filter model.ResourceQueryFilter) ([]model.Resources, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "GetResourcesAccessible called")
	// normalize filter (limit/offset) at usecase layer
	if filter.Limit <= 0 || filter.Limit > 200 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	groups, err := uc.GetAccessibleGroupUUIDs(c, userUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get accessible groups")
		return nil, err
	}
	if len(groups) == 0 {
		INFO(reqID, Mcode(SRNRSR1), "No accessible groups found")
		return []model.Resources{}, nil
	}
	filter.GroupUUIDs = groups
	resources, err := uc.resourceRepo.GetResources(c, filter)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get resources")
		return nil, err
	}
	INFO(reqID, Mcode(SRNRSR1), "GetResourcesAccessible succeeded")
	return resources, nil
}

// GetResourceAccessible returns a resource only if the user is a member of its group.
func (uc *resource) GetResourceAccessible(c *gin.Context, userUUID string, id string) (*model.Resources, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "GetResourceAccessible called")
	res, err := uc.resourceRepo.GetResource(c, id)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get resource")
		return nil, err
	}
	groups, err := uc.GetAccessibleGroupUUIDs(c, userUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get accessible groups")
		return nil, err
	}
	for _, g := range groups {
		if g == res.GroupUUID {
			INFO(reqID, Mcode(SRNRSR1), "GetResourceAccessible succeeded")
			return res, nil
		}
	}
	ERROR(reqID, Mcode(SRNRSR2), "Access denied")
	return nil, errors.New("access denied")
}

// CountResourcesAccessible counts resources in groups the user has access to.
func (uc *resource) CountResourcesAccessible(c *gin.Context, userUUID string, filter model.ResourceQueryFilter) (int64, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CountResourcesAccessible called")
	// normalize pagination/filter defaults in usecase
	if filter.Limit <= 0 || filter.Limit > 200 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	groups, err := uc.GetAccessibleGroupUUIDs(c, userUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get accessible groups")
		return 0, err
	}
	if len(groups) == 0 {
		INFO(reqID, Mcode(SRNRSR1), "No accessible groups found")
		return 0, nil
	}
	filter.GroupUUIDs = groups
	count, err := uc.resourceRepo.CountResources(c, filter)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to count resources")
		return 0, err
	}
	INFO(reqID, Mcode(SRNRSR1), "CountResourcesAccessible succeeded")
	return count, nil
}

// CreateResourceWithAccess allows creation when the user has editor or owner role on the group.
func (uc *resource) CreateResourceWithAccess(c *gin.Context, userUUID string, resource *model.Resources) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CreateResourceWithAccess called")
	if resource.GroupUUID == "" {
		ERROR(reqID, Mcode(SRNRSR2), "GroupUUID is required")
		return errors.New("group_uuid required")
	}
	member, err := uc.getMemberForUserAndGroup(c, userUUID, resource.GroupUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get member")
		return err
	}
	if member == nil {
		ERROR(reqID, Mcode(SRNRSR2), "Access denied - not a member")
		return errors.New("access denied")
	}
	if member.Role != "owner" && member.Role != "editor" {
		ERROR(reqID, Mcode(SRNRSR2), "Insufficient role to create resource")
		return errors.New("insufficient role to create resource")
	}
	err = uc.resourceRepo.CreateResource(c, resource)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to create resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "CreateResourceWithAccess succeeded")
	return nil
}

// UpdateResourceWithAccess allows update when the user has editor or owner role on the owning group.
func (uc *resource) UpdateResourceWithAccess(c *gin.Context, userUUID string, resource *model.Resources) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "UpdateResourceWithAccess called")
	// Ensure resource exists
	existing, err := uc.resourceRepo.GetResource(c, resource.UUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get existing resource")
		return err
	}
	member, err := uc.getMemberForUserAndGroup(c, userUUID, existing.GroupUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get member")
		return err
	}
	if member == nil {
		ERROR(reqID, Mcode(SRNRSR2), "Access denied - not a member")
		return errors.New("access denied")
	}
	if member.Role != "owner" && member.Role != "editor" {
		ERROR(reqID, Mcode(SRNRSR2), "Insufficient role to update resource")
		return errors.New("insufficient role to update resource")
	}
	err = uc.resourceRepo.UpdateResource(c, resource)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to update resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "UpdateResourceWithAccess succeeded")
	return nil
}

// DeleteResourceWithAccess allows delete when the user has editor or owner role on the owning group.
func (uc *resource) DeleteResourceWithAccess(c *gin.Context, userUUID string, id string) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "DeleteResourceWithAccess called")
	existing, err := uc.resourceRepo.GetResource(c, id)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get existing resource")
		return err
	}
	member, err := uc.getMemberForUserAndGroup(c, userUUID, existing.GroupUUID)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get member")
		return err
	}
	if member == nil {
		ERROR(reqID, Mcode(SRNRSR2), "Access denied - not a member")
		return errors.New("access denied")
	}
	if member.Role != "owner" && member.Role != "editor" {
		ERROR(reqID, Mcode(SRNRSR2), "Insufficient role to delete resource")
		return errors.New("insufficient role to delete resource")
	}
	err = uc.resourceRepo.DeleteResource(c, id)
	if err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to delete resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "DeleteResourceWithAccess succeeded")
	return nil
}

// Administrative passthrough methods (no access checks) used by private/admin controllers
func (uc *resource) GetResourcesAccessibleAdmin(c *gin.Context, _ string, filter model.ResourceQueryFilter) ([]model.Resources, error) {
	// alias to repository GetResources: ignoring user filtering
	return uc.resourceRepo.GetResources(c, filter)
}

func (uc *resource) GetResourceAccessibleAdmin(c *gin.Context, _ string, id string) (*model.Resources, error) {
	return uc.resourceRepo.GetResource(c, id)
}

func (uc *resource) CountResourcesAccessibleAdmin(c *gin.Context, _ string, filter model.ResourceQueryFilter) (int64, error) {
	return uc.resourceRepo.CountResources(c, filter)
}

func (uc *resource) CreateResourceWithAdmin(c *gin.Context, _ string, resource *model.Resources) error {
	return uc.resourceRepo.CreateResource(c, resource)
}

func (uc *resource) UpdateResourceWithAdmin(c *gin.Context, _ string, resource *model.Resources) error {
	return uc.resourceRepo.UpdateResource(c, resource)
}

func (uc *resource) DeleteResourceWithAdmin(c *gin.Context, _ string, id string) error {
	return uc.resourceRepo.DeleteResource(c, id)
}
