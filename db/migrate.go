package db

import (
	"embed"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var fs embed.FS

func Migrate(db *sqlx.DB) {
	goose.SetBaseFS(fs)

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Couldn't set goose dialect", "err", err)
		panic(err)
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		slog.Error("Applying database migrations failed", "err", err)
		panic(err)
	}
}
