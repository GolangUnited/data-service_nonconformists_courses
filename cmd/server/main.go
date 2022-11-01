package main

import (
	"context"
	"licheropew/golang-united-courses/pkg/api"
	"licheropew/golang-united-courses/pkg/courses"
	"licheropew/golang-united-courses/pkg/db"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go runRest()
	runGrpc()
}

func runRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterCoursesHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		panic(err)
	}
	log.Printf("server listening at 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
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
