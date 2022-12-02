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
	err = p.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return err
	}
	result := p.DB.Exec("SELECT 1 FROM pg_type WHERE typname = 'status';")
	switch {
	case result.RowsAffected == 0:
		err = p.DB.Exec("CREATE TYPE status AS ENUM ('unknown', 'joined', 'started', 'finished', 'declined');").Error
		if err != nil {
			return err
		}
	case result.Error != nil:
		return result.Error
	}
	err = p.DB.AutoMigrate(models.Course{}, models.UserCourse{})
	if err != nil {
		return err
	}
	return err
}

func (p *PostgreSql) Close() error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	db.Close()
	return nil
}
