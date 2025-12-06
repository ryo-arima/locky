package model

import "time"

type Roles struct {
	ID        uint   `gorm:"primaryKey,autoIncrement"`
	UUID      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
