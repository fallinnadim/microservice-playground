package outbound

import (
	"context"

	"github.com/fallinnadim/order-service/internal/adapter/outbound/order"
)

type OrderRepository interface {
	CreateOrder(context.Context, order.OrderInput) (string, error)
	UpdateOrder(context.Context, string, string) error
}
