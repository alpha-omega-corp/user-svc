package main

import (
	"github.com/alpha-omega-corp/authentication-svc/pkg/config"
	"github.com/alpha-omega-corp/authentication-svc/pkg/database"
	"github.com/alpha-omega-corp/authentication-svc/pkg/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	c, err := config.LoadConfig("dev")
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	dbHandler := database.NewHandler(c.DB)
	db := dbHandler.Database()

	defer func(database *bun.DB) {
		err := database.Close()
		if err != nil {
			log.Fatalln("Failed to close database:", err)
		}
	}(db)

	appCli := &cli.App{
		Name:  "authentication-svc",
		Usage: "bootstrap the service",
		Commands: []*cli.Command{
			migrateCommand(db),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func migrateCommand(db *bun.DB) *cli.Command {
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
					return db.ResetModel(c.Context, (*models.User)(nil))
				},
			},
		},
	}
}
