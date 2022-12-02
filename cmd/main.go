package main

import (
	"fmt"
	"golang-united-courses/config"
	"golang-united-courses/internal/api"
	"golang-united-courses/internal/database"
	"golang-united-courses/internal/interfaces"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	var myDb interfaces.CourseManager
	var dbUrl string
	//get APP configuration
	conf := config.New()
	//select db type
	switch conf.DBType {
	case "postgres":
		myDb = database.NewPgSql()
		dbUrl = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			conf.DBCfg.Host,
			conf.DBCfg.User,
			conf.DBCfg.Password,
			conf.DBCfg.Name,
			conf.DBCfg.Port,
			conf.DBCfg.SslMode,
			conf.DBCfg.Timezone,
		)
	default:
		log.Fatal("Database type not implemented")
	}
	// connect to DB
	if err := myDb.Init(dbUrl); err != nil {
		log.Printf("Database connection error: %s", err.Error())
	}
	defer myDb.Close()
	// create Course Server API
	courseServer := api.New(myDb)
	// reate and run GPRC-server with Course API
	grpcServer := grpc.NewServer()
	api.RegisterCoursesServer(grpcServer, courseServer)
	var port = "8080"
	if conf.Port != "" {
		port = conf.Port
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
