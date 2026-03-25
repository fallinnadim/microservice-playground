package outbound

import (
	"context"
	"time"

	"github.com/fallinnadim/order-service/internal/adapter/outbound/item"
	"github.com/fallinnadim/order-service/internal/domain"
)

type ItemRepository interface {
	FindByIds(ctx context.Context, ids []string) ([]*domain.Item, error)
}

type ItemCacheRepository interface {
	GetItems(ctx context.Context, keys []string) (map[string]item.ItemCache, error)
	SetItems(ctx context.Context, items map[string]item.ItemCache, ttl time.Duration) error
}
