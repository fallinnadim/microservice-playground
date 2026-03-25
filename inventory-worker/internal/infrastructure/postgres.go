package infrastructure

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(dsn string, log *slog.Logger) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error("failed to create DB pool", "error", err)
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Error("database unreachable or unhealthy", "error", err)
		os.Exit(1)
	}

	log.Info("✅ Database connected successfully")
	return pool
}
