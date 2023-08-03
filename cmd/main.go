package main

import (
	"github.com/alpha-omega-corp/authentication-svc/pkg/proto"
	"github.com/alpha-omega-corp/authentication-svc/pkg/services"
	"github.com/alpha-omega-corp/services/config"
	"github.com/alpha-omega-corp/services/database"
	"github.com/alpha-omega-corp/services/server"
	"google.golang.org/grpc"
)

func main() {
	c := config.Get("dev")

	err := server.NewGRPC(c.AUTH, c, func(h *database.Handler, grpc *grpc.Server) {
		s := services.NewServer(h.Database())
		proto.RegisterAuthServiceServer(grpc, s)
	})

	if err != nil {
		panic(err)
	}
}
