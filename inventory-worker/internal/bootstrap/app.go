package bootstrap

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/fallinnadim/inventory-worker/config"
	"github.com/fallinnadim/inventory-worker/internal/adapter/outbound/item"
	"github.com/fallinnadim/inventory-worker/internal/domain"
	"github.com/fallinnadim/inventory-worker/internal/infrastructure"
	"github.com/fallinnadim/inventory-worker/internal/logger"
	"github.com/fallinnadim/inventory-worker/internal/usecase"
	"github.com/segmentio/kafka-go"
)

type App struct {
	Logger       *slog.Logger
	Reader       *kafka.Reader
	WorkerNumber int
	ItemUsecase  *usecase.ItemUsecase
}

func NewApp(cfg *config.Config) (*App, error) {
	l := logger.New(logger.Config{
		Level: cfg.LogLevel,
	})
	db := infrastructure.NewPostgres(cfg.DbUrl, l)
	rdb := infrastructure.NewRedis(cfg.RedisUrl, l)
	itemRepo := item.NewItemRepository(db)
	cacheRepo := item.NewCacheRepository(rdb)
	itemUC := usecase.NewItemUsecase(itemRepo, cacheRepo)

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "order-worker-group",
		Topic:    "order.created",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &App{
		Logger:       l,
		Reader:       reader,
		WorkerNumber: cfg.WorkerNumber,
		ItemUsecase:  itemUC,
	}, nil
}

func (a *App) RunWorker(ctx context.Context) {
	const workerCount = 10
	var wg sync.WaitGroup

	a.Logger.Info("starting Kafka consumer pool",
		"topic", "order.created",
		"workers", workerCount,
	)

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			a.worker(ctx, workerID)
		}(i)
	}

	wg.Wait()
	a.Logger.Info("all workers have shut down")
}

func (a *App) worker(ctx context.Context, id int) {
	wLog := a.Logger.With("worker_id", id)
	wLog.Debug("worker started")

	for {
		msg, err := a.Reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				wLog.Debug("worker shutting down")
				return
			}
			wLog.Error("error reading message", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}
		var order []domain.OrderItem
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			wLog.Error("failed to unmarshal order", "error", err)
			continue
		}
		err = a.ItemUsecase.UpdateInventory(ctx, order)
		if err != nil {
			wLog.Error("business logic failure", "order_id", msg.Key, "error", err)
			// Decide here: continue (skip) or retry logic?
			continue
		}
		err = a.ItemUsecase.InvalidateCache(ctx, order)
		wLog.Info("successfully processed order", "order_id", msg.Key)
	}
}

func (a *App) Close() {
	if a.Reader != nil {
		if err := a.Reader.Close(); err != nil {
			a.Logger.Error("error closing kafka reader", "error", err)
		} else {
			a.Logger.Info("kafka reader closed successfully")
		}
	}
}
