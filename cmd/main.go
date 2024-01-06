package main

import (
	"github.com/alpha-omega-corp/auth-svc/pkg/models"
	"github.com/alpha-omega-corp/auth-svc/pkg/services"
	"github.com/alpha-omega-corp/auth-svc/pkg/utils"
	"github.com/alpha-omega-corp/auth-svc/proto"
	"github.com/alpha-omega-corp/services/database"
	"github.com/alpha-omega-corp/services/server"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {

	v := viper.New()
	cManager := server.NewConfigManager(v)

	config, err := cManager.HostsConfig()
	if err != nil {
		panic(err)
	}

	authConfig, err := cManager.AuthConfig()
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(config.Auth.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := server.NewGRPC(config.Auth.Host, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		s := services.NewServer(db, utils.NewAuthWrapper(authConfig.JwtSecret))
		proto.RegisterAuthServiceServer(grpc, s)
	}); err != nil {
		panic(err)
	}
}
