package pubsub

import (
	"context"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		slog.Error("Error connecting to Redis", "err", err)
		panic(err)
	}
	return client
}
