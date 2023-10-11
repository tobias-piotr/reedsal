package main

import (
	"os"
	"reedsal/cmd/server"
	"reedsal/db"
	"reedsal/pubsub"
)

func main() {
	conn := db.GetDB(os.Getenv("DATABASE_DSN"))
	db.Migrate(conn)
	redis := pubsub.NewRedisClient()
	server.NewServer().WithMiddleware().WithDB(conn).WithRedis(redis).Mount().Serve()
}
