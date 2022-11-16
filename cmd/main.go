package main

import (
	"golang-united-courses/config"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/interfaces.go"
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
	var myDb interfaces.CourseManager
	mdb := new(courses.PostgreSql)
	mdb.Init(url)
	myDb = mdb
	defer myDb.Close()
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
