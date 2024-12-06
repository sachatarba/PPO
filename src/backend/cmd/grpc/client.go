package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/sachatarba/course-db/internal/api/grpc"
	"github.com/sachatarba/course-db/internal/config"
	postrgres_adapter "github.com/sachatarba/course-db/internal/postrgres"
	"github.com/sachatarba/course-db/internal/repository"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	server "github.com/sachatarba/course-db/internal/grpc"
	"github.com/sachatarba/course-db/internal/service"
)

func main() {
	conf := config.NewConfFromEnv()

	postgresConnector := postrgres_adapter.PostgresConnector{
		Conf: conf.PostgresConf,
	}

	// db, err := postgresConnector.Connect()
	// if err != nil {
	// 	log.Print("Cant connect postgres", err)
	// 	return
	// }
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = postgresConnector.Connect()
		if err != nil {
			log.Print("Cant connect postgres", err)
			time.Sleep(time.Second)
			// return
		}
	}
	if err != nil {
		log.Print("Cant connect postgres", err)
		return
	}

	clientRepo := repository.NewClientRepo(db)
	clientService := service.NewClientService(clientRepo)

	grpcServer := grpc.NewServer()
	pb.RegisterClientServiceServer(grpcServer, server.NewGRPCClientServer(clientService))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s",
		conf.GrpcClientsServerConfig.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server is running on port :", conf.GrpcClientsServerConfig.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
