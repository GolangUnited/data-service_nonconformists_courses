package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Name     string
	Password string
	Port     string
	User     string
	Timezone string
	SslMode  string
}

type Config struct {
	DBCfg  DBConfig
	DBType string
	Port   string
}

func New() *Config {
	return &Config{
		DBCfg: DBConfig{
			Host:     os.Getenv("COURSES_DB_HOST"),
			Name:     os.Getenv("COURSES_DB_NAME"),
			Password: os.Getenv("COURSES_DB_PASSWORD"),
			Port:     os.Getenv("COURSES_DB_PORT"),
			User:     os.Getenv("COURSES_DB_USER"),
			Timezone: os.Getenv("COURSES_DB_TZ"),
			SslMode:  os.Getenv("COURSES_DB_SSLMODE"),
		},
		DBType: os.Getenv("COURSES_DB_TYPE"),
		Port:   os.Getenv("COURSES_PORT"),
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}
