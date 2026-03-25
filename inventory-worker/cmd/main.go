package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/fallinnadim/inventory-worker/config"
	"github.com/fallinnadim/inventory-worker/internal/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") == "production" {
		_ = godotenv.Load("inventory-worker/.env")
	} else {
		_ = godotenv.Load("inventory-worker/.env.local")
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	app, err := bootstrap.NewApp(cfg)
	if err != nil {
		slog.Error("failed to init app", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		app.Logger.Info("Cleaning up resources...")
		app.Close()
	}()

	app.RunWorker(ctx)

	app.Logger.Info("Shutdown complete.")
}
