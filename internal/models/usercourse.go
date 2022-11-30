package models

import (
	"time"

	"github.com/google/uuid"
)

type UserCourse struct {
	CourseID        uuid.UUID `gorm:"primaryKey"`
	UserID          uuid.UUID `gorm:"primaryKey"`
	CreatedAt       time.Time
	StartDate       time.Time
	FinishDate      time.Time
	PercentFinished uint32
	Status          string
}
