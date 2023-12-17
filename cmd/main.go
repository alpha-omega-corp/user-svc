package main

import (
	"github.com/alpha-omega-corp/auth-svc/pkg/config"
	"github.com/alpha-omega-corp/auth-svc/pkg/models"
	"github.com/alpha-omega-corp/auth-svc/pkg/services"
	"github.com/alpha-omega-corp/auth-svc/pkg/utils"
	"github.com/alpha-omega-corp/auth-svc/proto"
	"github.com/alpha-omega-corp/services/database"
	"github.com/alpha-omega-corp/services/server"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(c.DSN)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
		(*models.RoleToPermission)(nil),
	)

	if err := server.NewGRPC(c.HOST, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		s := services.NewServer(db, utils.NewAuthWrapper(c.KEY))
		proto.RegisterAuthServiceServer(grpc, s)
	}); err != nil {
		panic(err)
	}
}
