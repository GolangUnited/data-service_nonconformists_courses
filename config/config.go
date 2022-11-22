package config

import (
	"os"
)

type DBConfig struct {
	COURSES_DB_HOST     string
	COURSES_DB_NAME     string
	COURSES_DB_PASSWORD string
	COURSES_DB_PORT     string
	COURSES_DB_USER     string
	COURSES_DB_TZ       string
	COURSES_DB_SSLMODE  string
}

type Config struct {
	DBCfg  DBConfig
	DBType string
}

var cfg *Config

func New() *Config {
	if cfg == nil {
		cfg = &Config{
			DBCfg: DBConfig{
				COURSES_DB_HOST:     os.Getenv("COURSES_DB_HOST"),
				COURSES_DB_NAME:     os.Getenv("COURSES_DB_NAME"),
				COURSES_DB_PASSWORD: os.Getenv("COURSES_DB_PASSWORD"),
				COURSES_DB_PORT:     os.Getenv("COURSES_DB_PORT"),
				COURSES_DB_USER:     os.Getenv("COURSES_DB_USER"),
				COURSES_DB_SSLMODE:  os.Getenv("COURSES_DB_SSLMODE"),
				COURSES_DB_TZ:       os.Getenv("COURSES_DB_TZ"),
			},
			DBType: os.Getenv("COURSES_DB_TYPE"),
		}
	}
	return cfg
}
