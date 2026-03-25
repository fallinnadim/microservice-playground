package bootstrap

import (
	"github.com/fallinnadim/payment-service/config"
	httpAdapter "github.com/fallinnadim/payment-service/internal/adapter/inbound/http"
	"github.com/fallinnadim/payment-service/internal/logger"
	"github.com/go-playground/validator/v10"
)

type App struct {
	Router *httpAdapter.Handler
}

func NewApp(cfg *config.Config) (*App, error) {
	log := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})
	validator := validator.New()
	config.InitOpentel()

	handler := httpAdapter.NewHandler(log, validator)

	return &App{
		Router: handler,
	}, nil
}
