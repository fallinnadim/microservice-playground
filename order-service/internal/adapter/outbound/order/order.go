package order

import "github.com/fallinnadim/order-service/internal/domain"

type OrderInput struct {
	UserId string
	Status string
	Items  []domain.OrderItem
}
