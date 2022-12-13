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
			Host:     getEnv("COURSES_DB_HOST", ""),
			Name:     getEnv("COURSES_DB_NAME", ""),
			Password: getEnv("COURSES_DB_PASSWORD", ""),
			Port:     getEnv("COURSES_DB_PORT", ""),
			User:     getEnv("COURSES_DB_USER", ""),
			Timezone: getEnv("COURSES_DB_TZ", ""),
			SslMode:  getEnv("COURSES_DB_SSLMODE", ""),
		},
		DBType: getEnv("COURSES_DB_TYPE", ""),
		Port:   getEnv("COURSES_PORT", "8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}
