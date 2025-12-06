package model

import "time"

// Resources represents a resource entity that belongs to a group.
// Access to resources is determined by the user's membership in the owning group.
type Resources struct {
	ID        uint   `gorm:"primaryKey,autoIncrement"`
	UUID      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"uniqueIndex;not null"`
	Type      string
	GroupUUID string `gorm:"index;not null"` // Reference to the owning group
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// ResourceQueryFilter defines filter conditions for resource queries.
// This struct lives in the model package so it can be shared between
// controller/usecase/repository layers without duplicating definitions.
type ResourceQueryFilter struct {
	ID         *uint
	UUID       *string
	Name       *string
	NamePrefix *string
	NameLike   *string
	Type       *string
	GroupUUID  *string
	GroupUUIDs []string // Multiple group UUIDs for access control
	Limit      int
	Offset     int
}
