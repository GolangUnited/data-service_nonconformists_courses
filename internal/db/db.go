package db

import (
	"golang-united-courses/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(dsn string) (*Handler, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Course{})
	if err != nil {
		return nil, err
	}
	return &Handler{db}, err
}

func (h *Handler) Close() error {
	db, err := h.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
