package outbound

import (
	"context"

	"github.com/fallinnadim/order-service/internal/adapter/outbound/order"
)

type PaymentService interface {
	Pay(ctx context.Context, req order.PaymentRequest) (*order.PaymentResponse, error)
}
