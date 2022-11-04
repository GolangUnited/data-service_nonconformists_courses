package config

import (
	"fmt"
	"os"
)

func GetConfig() string {
	dbHost := os.Getenv("COURSES_DB_HOST")
	dbUser := os.Getenv("COURSES_DB_USER")
	dbPassword := os.Getenv("COURSES_DB_PASSWORD")
	dbName := os.Getenv("COURSES_DB_NAME")
	dbPort := os.Getenv("COURSES_DB_PORT")
	// database URL
	url := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s  sslmode=disable TimeZone=Europe/Moscow", dbHost, dbUser, dbPassword, dbName, dbPort)
	return url
}
