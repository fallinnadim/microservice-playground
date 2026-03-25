package http

import (
	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/handler"
	"github.com/fallinnadim/order-service/internal/usecase"
)

type Handler struct {
	auth        *handler.AuthHandler
	order       *handler.OrderHandler
	rateLimitUC *usecase.RateLimitUsecase
}

func NewHandler(
	auth *handler.AuthHandler,
	order *handler.OrderHandler,
	rateLimitUC *usecase.RateLimitUsecase,
) *Handler {
	return &Handler{
		auth, order, rateLimitUC,
	}
}
