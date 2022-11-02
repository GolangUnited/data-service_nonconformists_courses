package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title       string
	Description string
	CreatedBy   string
	UpdatedBy   string
	DeletedBy   string
}
