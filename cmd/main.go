package main

import (
	"github.com/alpha-omega-corp/authentication-svc/pkg/config"
	"github.com/alpha-omega-corp/authentication-svc/pkg/services"
	"github.com/alpha-omega-corp/authentication-svc/proto"
	"github.com/alpha-omega-corp/services/server"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	h := config.NewConnection(c.DSN)

	if err := server.NewGRPC(c.HOST, h.Database(), func(db *bun.DB, grpc *grpc.Server) {
		s := services.NewServer(db)
		proto.RegisterAuthServiceServer(grpc, s)
	}); err != nil {
		panic(err)
	}
}
