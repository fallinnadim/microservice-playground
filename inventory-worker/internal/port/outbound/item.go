package outbound

import (
	"context"

	"github.com/fallinnadim/inventory-worker/internal/domain"
)

type ItemRepository interface {
	UpdateStocks(ctx context.Context, items []domain.OrderItem) error
}
