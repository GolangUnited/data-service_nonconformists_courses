package main

import (
	"golang-united-courses/config"
	"golang-united-courses/pkg/api"
	"golang-united-courses/pkg/courses"
	"golang-united-courses/pkg/db"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	runGrpc()
}

func runGrpc() {
	url := config.GetConfig()
	myDb, err := db.Init(url)
	if err != nil {
		log.Fatal(err)
	}
	defer myDb.Close()
	myCourseServer := &courses.CourseServer{
		Course: myDb,
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
