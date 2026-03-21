package bootstrap

import (
	"github.com/fallinnadim/order-service/config"
	httpAdapter "github.com/fallinnadim/order-service/internal/adapter/inbound/http"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/auth"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/auth/jwt"
	ratelimit "github.com/fallinnadim/order-service/internal/adapter/outbound/rate_limit"
	"github.com/fallinnadim/order-service/internal/infrastructure"
	"github.com/fallinnadim/order-service/internal/logger"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	rateLimitAdapter := ratelimit.NewRateLimitAdapter(redis)
	userRepo := auth.NewUserRepository(db)
	authAdapter := jwt.NewJWTAuthAdapter(cfg.JWTSecret, cfg.JWTDuration)
	authUC := usecase.NewAuthUsecase(authAdapter, userRepo)
	rateLimitRefillRate := float64(cfg.RateLimitPerMinute) / 60.0
	rateLimitUC := usecase.NewRateLimitUsecase(
		rateLimitAdapter,
		cfg.RateLimitCapacity,
		rateLimitRefillRate,
	)
	pingUC := usecase.NewPingUsecase(db)
	handler := httpAdapter.NewHandler(
		pingUC,
		authUC,
		rateLimitUC,
		log,
	)

	return &App{
		Router: handler,
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
