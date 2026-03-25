package bootstrap

import (
	"net/http"
	"strings"
	"time"

	"github.com/fallinnadim/order-service/config"
	httpAdapter "github.com/fallinnadim/order-service/internal/adapter/inbound/http"
	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/handler"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/auth"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/auth/argon2"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/auth/jwt"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/item"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/kafka"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/order"
	ratelimit "github.com/fallinnadim/order-service/internal/adapter/outbound/rate_limit"
	"github.com/fallinnadim/order-service/internal/infrastructure"
	"github.com/fallinnadim/order-service/internal/logger"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type App struct {
	Router *httpAdapter.Handler
	Db     *pgxpool.Pool
	Redis  *redis.Client
}

func NewApp(cfg *config.Config) (*App, error) {
	log := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})
	db := infrastructure.NewPostgres(cfg.DbUrl, log)
	redis := infrastructure.NewRedis(cfg.RedisUrl, log)
	config.InitOpentel()
	tlsConfig, _ := config.LoadTLSConfig(cfg)
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: otelhttp.NewTransport(&http.Transport{
			TLSClientConfig: tlsConfig,
		}),
	}

	rateLimitAdapter := ratelimit.NewRateLimitAdapter(redis)
	paymentAdapter := order.NewPaymentAdapter(cfg.PaymentServiceURL, httpClient)
	userRepo := auth.NewUserRepository(db)
	itemRepo := item.NewItemRepository(db)
	orderRepo := order.NewOrderRepository(db)
	itemCacheRepo := item.NewItemCacheRepository(redis)
	jwtAdapter := jwt.NewJWTAuthAdapter(cfg.JWTSecret, cfg.JWTDuration)
	argon2Adapter := argon2.NewJWTAuthAdapter()
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}
	kafkaProducer := kafka.NewKafkaProducer(brokers)

	authUC := usecase.NewAuthUsecase(jwtAdapter, argon2Adapter, userRepo)
	rateLimitRefillRate := float64(cfg.RateLimitPerMinute) / 60.0
	rateLimitUC := usecase.NewRateLimitUsecase(
		rateLimitAdapter,
		cfg.RateLimitCapacity,
		rateLimitRefillRate,
	)
	orderUC := usecase.NewOrderUsecase(itemRepo, itemCacheRepo, orderRepo, paymentAdapter, kafkaProducer)

	validate := validator.New()
	authHandler := handler.NewAuthHandler(validate, authUC, log)
	orderHandler := handler.NewOrderHandler(validate, orderUC, log)

	topLevelHandler := httpAdapter.NewHandler(authHandler, orderHandler, rateLimitUC)

	return &App{
		Router: topLevelHandler,
		Db:     db,
		Redis:  redis,
	}, nil
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
	if a.Redis != nil {
		a.Redis.Close()
	}
}
