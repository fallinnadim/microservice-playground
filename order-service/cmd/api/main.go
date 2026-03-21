package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fallinnadim/order-service/config"
	httpAdapter "github.com/fallinnadim/order-service/internal/adapter/inbound/http"
	"github.com/fallinnadim/order-service/internal/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type Something struct {
	Message string `json:"message"`
}

func main() {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load("order-service/.env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(cfg.GinMode)
	app, _ := bootstrap.NewApp(cfg)
	defer app.Close()
	router := httpAdapter.NewRouter(app.Router)

	router.Run(fmt.Sprintf(":%d", cfg.Port))
}
