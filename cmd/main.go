package main

import (
	"licheropew/golang-united-courses/internal/api"
	"licheropew/golang-united-courses/internal/courses"
	"licheropew/golang-united-courses/internal/db"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	runGrpc()
}

func runGrpc() {
	s := grpc.NewServer()
	h, err := db.Init("postgres://postgres:postgrespw@localhost:49155")
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()
	myCourse := &courses.Server{
		C: h,
	}

	api.RegisterCoursesServer(s, myCourse)

	c, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(c); err != nil {
		log.Fatal(err)
	}
}
