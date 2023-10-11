package db

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDB(dsn string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("Couldn't connect to the database", "err", err)
		panic(err)
	}
	return db
}
