package inbound

import "context"

type OrderUsecase interface {
	Order(context.Context, OrderInput) (string, error)
}

type OrderInput struct {
	UserId string
	Items  []ItemInput
}

type ItemInput struct {
	Id      string
	Ammount int
}
