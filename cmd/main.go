package main

import (
	localConfig "github.com/alpha-omega-corp/auth-svc/config"
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
	hostsConfig, err := localConfig.HostsConfig()
	if err != nil {
		panic(err)
	}

	authConfig, err := localConfig.AuthConfig()
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(hostsConfig.Auth.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := server.NewGRPC(hostsConfig.Auth.Host, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		s := services.NewServer(db, utils.NewAuthWrapper(authConfig.JwtSecret))
		proto.RegisterAuthServiceServer(grpc, s)
	}); err != nil {
		panic(err)
	}
}
