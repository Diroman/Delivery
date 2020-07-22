package main

import (
	"campus-delivery/config"
	"campus-delivery/domain/database"
	usecase2 "campus-delivery/domain/usecase"
	"fmt"
	"log"
	"os"

	"campus-delivery/clientApi/controller"
	"campus-delivery/clientApi/grpc"
)

func main() {
	conf := config.Read()
	server := configureServer(conf)
	if err := server.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		log.Printf("%v\n", err)
		os.Exit(1)
	}
}

func configureServer(conf config.Config) *grpc.GrpcServer {
	dbConf := conf.Database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.DBName)
	log.Printf("Connect to database {host: %v, port: %v} ...", dbConf.Host, dbConf.Port)
	db := database.NewDBClient(psqlInfo)
	if err := db.Connect(); err != nil {
		_ = fmt.Errorf("Database connection error: %v", err)
	}
	go db.DeleteUserByTimer()
	log.Printf("Success connect")

	usecase := usecase2.NewUseCase(db)
	newController := controller.NewController(usecase)
	server := grpc.NewServer(conf.Server.GrpcPort, newController)
	return server
}
