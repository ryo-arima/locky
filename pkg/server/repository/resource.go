package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/server/share"
	"gorm.io/gorm"
)

// ResourceQueryFilter moved to `pkg/entity/model/resource.go` as `model.ResourceQueryFilter`.
type Resource interface {
	GetResources(c *gin.Context, filter model.ResourceQueryFilter) ([]model.Resources, error)
	GetResource(c *gin.Context, id string) (*model.Resources, error)
	CountResources(c *gin.Context, filter model.ResourceQueryFilter) (int64, error)
	CreateResource(c *gin.Context, resource *model.Resources) error
	UpdateResource(c *gin.Context, resource *model.Resources) error
	DeleteResource(c *gin.Context, id string) error
}

type resource struct {
	db     *gorm.DB
	logger global.MCode
}

func NewResource(conf config.BaseConfig) Resource {
	return &resource{
		db:     conf.DBConnection,
		logger: global.SRNRSR1,
	}
}

func (r *resource) GetResources(c *gin.Context, filter model.ResourceQueryFilter) ([]model.Resources, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "GetResources called")
	q := r.db.Model(&model.Resources{}).Where("deleted_at IS NULL")

	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.Name != nil {
		q = q.Where("name = ?", *filter.Name)
	}
	if filter.NamePrefix != nil {
		q = q.Where("name LIKE ?", *filter.NamePrefix+"%")
	}
	if filter.NameLike != nil {
		q = q.Where("name LIKE ?", "%"+*filter.NameLike+"%")
	}
	if filter.Type != nil {
		q = q.Where("type = ?", *filter.Type)
	}
	if filter.GroupUUID != nil {
		q = q.Where("group_uuid = ?", *filter.GroupUUID)
	}
	if len(filter.GroupUUIDs) > 0 {
		q = q.Where("group_uuid IN ?", filter.GroupUUIDs)
	}

	q = q.Limit(filter.Limit).Offset(filter.Offset)

	var resources []model.Resources
	if err := q.Find(&resources).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get resources")
		return nil, err
	}
	INFO(reqID, Mcode(SRNRSR1), "GetResources succeeded")
	return resources, nil
}

func (r *resource) GetResource(c *gin.Context, id string) (*model.Resources, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "GetResource called")
	var resource model.Resources
	if err := r.db.Where("uuid = ? AND deleted_at IS NULL", id).First(&resource).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to get resource")
		return nil, err
	}
	INFO(reqID, Mcode(SRNRSR1), "GetResource succeeded")
	return &resource, nil
}

func (r *resource) CountResources(c *gin.Context, filter model.ResourceQueryFilter) (int64, error) {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CountResources called")
	q := r.db.Model(&model.Resources{}).Where("deleted_at IS NULL")

	if filter.ID != nil {
		q = q.Where("id = ?", *filter.ID)
	}
	if filter.UUID != nil {
		q = q.Where("uuid = ?", *filter.UUID)
	}
	if filter.Name != nil {
		q = q.Where("name = ?", *filter.Name)
	}
	if filter.NamePrefix != nil {
		q = q.Where("name LIKE ?", *filter.NamePrefix+"%")
	}
	if filter.NameLike != nil {
		q = q.Where("name LIKE ?", "%"+*filter.NameLike+"%")
	}
	if filter.Type != nil {
		q = q.Where("type = ?", *filter.Type)
	}
	if filter.GroupUUID != nil {
		q = q.Where("group_uuid = ?", *filter.GroupUUID)
	}
	if len(filter.GroupUUIDs) > 0 {
		q = q.Where("group_uuid IN ?", filter.GroupUUIDs)
	}

	var count int64
	if err := q.Count(&count).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to count resources")
		return 0, err
	}
	INFO(reqID, Mcode(SRNRSR1), "CountResources succeeded")
	return count, nil
}

func (r *resource) CreateResource(c *gin.Context, resource *model.Resources) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "CreateResource called")
	if err := r.db.Create(resource).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to create resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "CreateResource succeeded")
	return nil
}

func (r *resource) UpdateResource(c *gin.Context, resource *model.Resources) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "UpdateResource called")
	if err := r.db.Where("uuid = ?", resource.UUID).Updates(resource).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to update resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "UpdateResource succeeded")
	return nil
}

func (r *resource) DeleteResource(c *gin.Context, id string) error {
	reqID := share.GetRequestID(c)
	INFO(reqID, Mcode(SRNRSR1), "DeleteResource called")
	if err := r.db.Model(&model.Resources{}).Where("uuid = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		ERROR(reqID, Mcode(SRNRSR2), "Failed to delete resource")
		return err
	}
	INFO(reqID, Mcode(SRNRSR1), "DeleteResource succeeded")
	return nil
}
