package main

import (
	"golang-united-courses/config"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/repositories/courses"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	runGrpc()
}

func runGrpc() {
	url := config.GetConfig()
	myDb := courses.PostgreSql{}
	err := myDb.Init(url)
	if err != nil {
		log.Fatal(err)
	}
	defer myDb.Close()
	myCourseServer := &api.CourseServer{
		DB: &myDb,
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
