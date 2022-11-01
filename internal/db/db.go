package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) (*Handler, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Handler{db}, err
}

func (h *Handler) AutoMigrate(i ...interface{}) error {
	return h.DB.AutoMigrate(i)
}

func (h *Handler) Close() error {
	db, err := h.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
