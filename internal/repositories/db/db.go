package db

import (
	"golang-united-courses/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSql struct {
	DB *gorm.DB
}

func (p *PostgreSql) Init(dsn string) error {
	var err error
	p.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err = p.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}
	result := p.DB.Exec("SELECT 1 FROM pg_type WHERE typname = 'status';")
	switch {
	case result.RowsAffected == 0:
		if err = p.DB.Exec("CREATE TYPE status AS ENUM ('unknown', 'joined', 'started', 'finished', 'declined');").Error; err != nil {
			return err
		}
	case result.Error != nil:
		return result.Error
	}
	if err = p.DB.AutoMigrate(models.Course{}, models.UserCourse{}); err != nil {
		return err
	}
	return nil
}

func (p *PostgreSql) Close() error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
