package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fallinnadim/payment-service/config"
	httpAdapter "github.com/fallinnadim/payment-service/internal/adapter/inbound/http"
	"github.com/fallinnadim/payment-service/internal/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") == "production" {
		_ = godotenv.Load("payment-service/.env")
	} else {
		_ = godotenv.Load("payment-service/.env.local")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(cfg.GinMode)
	app, _ := bootstrap.NewApp(cfg)
	router := httpAdapter.NewRouter(app.Router)
	tlsConfig, _ := config.LoadTLSConfig(cfg)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", cfg.Port),
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Printf("payment-service running with mTLS on :%d", cfg.Port)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
