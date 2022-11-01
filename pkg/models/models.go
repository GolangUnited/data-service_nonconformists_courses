package models

import "time"

type Course struct {
	Id          uint64 `json:"id" gorm:"primaryKey"`
	Title       string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}
