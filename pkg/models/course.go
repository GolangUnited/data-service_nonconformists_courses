package models

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID          string `gorm:"primarykey"`
	UserID      string `gorm:"index"`
	Title       string
	Description string
	CreatedBy   string
	UpdatedBy   string
	DeletedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
