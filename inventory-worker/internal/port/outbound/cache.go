package outbound

import (
	"context"

	"github.com/fallinnadim/inventory-worker/internal/domain"
)

type CacheRepository interface {
	InvalidateCache(ctx context.Context, items []domain.OrderItem) error
}
