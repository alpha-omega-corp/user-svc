package main

import (
	"fmt"
	"github.com/alpha-omega-corp/authentication-svc/pkg/config"
	"github.com/alpha-omega-corp/authentication-svc/pkg/database"
	"github.com/alpha-omega-corp/authentication-svc/pkg/services"
	"github.com/alpha-omega-corp/authentication-svc/pkg/services/pb"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	c, err := config.LoadConfig("dev")
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	listen, err := net.Listen("tcp", c.HOST)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	db := database.NewHandler(c.DB).Database()

	server := services.NewServer(db)
	pb.RegisterAuthServiceServer(grpcServer, server)

	fmt.Printf("Starting server... %v", c.HOST)
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("Failed to close db connection:", err)
		}
	}(db)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
