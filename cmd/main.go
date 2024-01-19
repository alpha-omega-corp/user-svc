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
	cManager := config.NewHandler()
	cHost, err := cManager.Manager().Hosts()
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(cHost.User.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := svc.NewGRPC(cHost.User.Host, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		cUser, err := cManager.Manager().UserService()
		if err != nil {
			panic(err)
		}

		userServer := server.NewServer(db, utils.NewAuthWrapper(cUser.JwtSecret))
		proto.RegisterUserServiceServer(grpc, userServer)
	}); err != nil {
		panic(err)
	}
}
