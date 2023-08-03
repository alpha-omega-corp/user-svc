package database

import (
	"database/sql"
	"github.com/alpha-omega-corp/authentication-svc/pkg/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"sync"
)

type Handler struct {
	dbOnce sync.Once
	db     *bun.DB

	config config.DbConfig
}

func NewHandler(c config.DbConfig) *Handler {
	return &Handler{
		config: c,
	}
}

func (h *Handler) Database() *bun.DB {
	h.dbOnce.Do(func() {
		dbConf := h.config
		driverOptions := pgdriver.NewConnector(
			pgdriver.WithAddr(dbConf.ADDR),
			pgdriver.WithDatabase(dbConf.NAME),
			pgdriver.WithUser(dbConf.USER),
			pgdriver.WithPassword(dbConf.PASS),
			pgdriver.WithTLSConfig(nil))

		conn := sql.OpenDB(driverOptions)
		db := bun.NewDB(conn, pgdialect.New())

		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithEnabled(true),
			bundebug.WithVerbose(true),
		))

		h.db = db
	})

	return h.db
}
