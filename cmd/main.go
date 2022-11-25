package main

import (
	"fmt"
	"golang-united-courses/config"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/interfaces.go"
	"golang-united-courses/internal/repositories/courses"
	"log"
	"net"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

func main() {
	runApp()
}

func runApp() {
	//get APP configuration
	conf := config.New()
	var myDb interfaces.CourseManager
	switch conf.DBType {
	case "postgres":
		myDb = new(courses.PostgreSql)
	default:
		log.Fatal("Database type not implemented")
	}
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		conf.DBCfg.Host,
		conf.DBCfg.User,
		conf.DBCfg.Password,
		conf.DBCfg.Name,
		conf.DBCfg.Port,
		conf.DBCfg.SslMode,
		conf.DBCfg.Timezone,
	)
	// connect to DB
	if err := myDb.Init(dbUrl); err != nil {
		log.Printf("Database connection error: %s", err.Error())
	}
	defer myDb.Close()
	// run GPRC-server
	myCourseServer := &api.CourseServer{
		DB: myDb,
	}
	s := grpc.NewServer()
	api.RegisterCoursesServer(s, myCourseServer)
	c, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(c); err != nil {
		log.Fatal(err)
	}
}
