package courses

import (
	"golang-united-courses/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
