package main

import (
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/courses"
	"golang-united-courses/internal/db"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	runGrpc()
}

func runGrpc() {
	db_host := os.Getenv("COURSES_DB_HOST")
	db_port := os.Getenv("COURSES_DB_PORT")
	db_user := os.Getenv("COURSES_DB_USER")
	db_password := os.Getenv("COURSES_DB_PASSWORD")
	db_name := os.Getenv("COURSES_DB_NAME")

	url := "host=" + db_host + " port=" + db_port + " user=" + db_user + " password=" + db_password + " database=" + db_name + " sslmode=disable TimeZone=Europe/Moscow"

	h, err := db.Init(url)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()
	myCourse := &courses.Server{
		C: h,
	}

	s := grpc.NewServer()

	api.RegisterCoursesServer(s, myCourse)

	c, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(c); err != nil {
		log.Fatal(err)
	}
}
