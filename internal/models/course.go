package models

import (
	"time"

	"github.com/google/uuid"
)

type Database struct {
	Course
	UserCourse
}

type Course struct {
	ID          uuid.UUID `gorm:"primarykey;type:uuid;default:public.uuid_generate_v4()"`
	Title       string
	Description string
	CreatedAt   time.Time
	IsDeleted   uint32
}

type UserCourse struct {
	CourseID        uuid.UUID `gorm:"primaryKey"`
	UserID          uuid.UUID `gorm:"primaryKey"`
	CreatedAt       time.Time
	IsDeleted       uint32
	StartDate       time.Time
	PercentFinished uint32
	FinishDate      time.Time
}
