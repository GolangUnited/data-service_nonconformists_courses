package db

import (
	"golang-united-courses/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func Init(dsn string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Course{})
	if err != nil {
		return nil, err
	}
	return &Database{db}, err
}

func (h *Database) Close() error {
	db, err := h.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
