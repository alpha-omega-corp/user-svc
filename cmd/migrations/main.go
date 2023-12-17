package main

import (
	"github.com/alpha-omega-corp/auth-svc/pkg/config"
	"github.com/alpha-omega-corp/auth-svc/pkg/models"
	"github.com/alpha-omega-corp/auth-svc/proto"
	"github.com/alpha-omega-corp/services/database"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	c, _ := config.LoadConfig()
	dbHandler := database.NewHandler(c.DSN)

	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dbHandler.Database())

	appCli := &cli.App{
		Name:  "authentication-svc",
		Usage: "bootstrap the service",
		Commands: []*cli.Command{
			migrateCommand(dbHandler.Database()),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func migrateCommand(db *bun.DB) *cli.Command {
	db.RegisterModel(
		(*models.UserToRole)(nil),
		(*models.RoleToPermission)(nil),
	)

	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrate.NewMigrations())
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "reset",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := db.ResetModel(c.Context,
						(*models.User)(nil),
						(*models.Role)(nil),
						(*models.UserToRole)(nil),
						(*models.Permission)(nil),
						(*models.RoleToPermission)(nil),
						(*models.Service)(nil),
					); err != nil {
						return err
					}

					svcMap := proto.ServicesEnum_value
					for name := range svcMap {
						_, err := db.NewInsert().Model(&models.Service{
							Name: name,
						}).Exec(c.Context)

						if err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}
}
