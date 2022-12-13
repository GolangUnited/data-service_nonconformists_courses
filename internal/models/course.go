package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID `gorm:"primarykey;type:uuid;default:public.uuid_generate_v4()"`
	Title       string
	Description string
	CreatedAt   time.Time
	IsDeleted   int32
}
