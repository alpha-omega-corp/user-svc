package main

import (
	"github.com/alpha-omega-corp/services/config"
	"github.com/alpha-omega-corp/services/database"
	svc "github.com/alpha-omega-corp/services/server"
	"github.com/alpha-omega-corp/user-svc/pkg/models"
	"github.com/alpha-omega-corp/user-svc/pkg/server"
	"github.com/alpha-omega-corp/user-svc/pkg/utils"
	"github.com/alpha-omega-corp/user-svc/proto"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {
	cHandler := config.NewHandler()
	env, err := cHandler.Environment("user")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(env.Host.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := svc.NewGRPC(env.Host.Url, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		auth := utils.NewAuthWrapper(env.Config.Viper.GetString("secret"))
		proto.RegisterUserServiceServer(grpc, server.NewServer(db, auth))
	}); err != nil {
		panic(err)
	}
}
