package model

import "time"

type Resources struct {
	ID        uint   `gorm:"primaryKey,autoIncrement"`
	UUID      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"uniqueIndex;not null"`
	Type      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
