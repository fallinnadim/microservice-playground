package infrastructure

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(dsn string, log *slog.Logger) *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		log.Error("redis unreachable or unhealthy", "error", err)
		os.Exit(1)
	}

	log.Info("✅ Redis connected successfully")
	return client
}
