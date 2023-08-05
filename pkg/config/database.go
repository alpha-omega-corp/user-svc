package config

import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"sync"
)

type Handler struct {
	dbOnce sync.Once
	db     *bun.DB

	dsn string
}

func NewConnection(dsn string) *Handler {
	return &Handler{
		dsn: dsn,
	}
}

func (h *Handler) Database() *bun.DB {
	h.dbOnce.Do(func() {
		driverOptions := pgdriver.NewConnector(
			pgdriver.WithDSN(h.dsn),
			pgdriver.WithTLSConfig(nil),
		)

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
