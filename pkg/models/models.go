package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	UserID      uint32
	Title       string
	Description string
	CreatedBy   string
	UpdatedBy   string
	DeletedBy   string
}
