package http

import (
	"log/slog"
	"net/http"

	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/port/outbound"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	pingUC      inbound.PingUsecase
	authUC      outbound.AuthService
	rateLimitUC *usecase.RateLimitUsecase
	log         *slog.Logger
}

func NewHandler(pingUC inbound.PingUsecase, authUC outbound.AuthService, rateLimitUC *usecase.RateLimitUsecase, log *slog.Logger) *Handler {
	return &Handler{
		pingUC, authUC, rateLimitUC, log,
	}
}

func (h *Handler) Ping(c *gin.Context) {
	h.log.Info("ping called")
	c.JSON(http.StatusOK, gin.H{
		"message": h.pingUC.Ping(),
	})
}
