package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/fallinnadim/payment-service/internal/adapter/inbound/http/request"
	"github.com/fallinnadim/payment-service/internal/adapter/inbound/http/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Log      *slog.Logger
	validate *validator.Validate
}

func NewHandler(
	log *slog.Logger,
	validate *validator.Validate,
) *Handler {
	return &Handler{log, validate}
}

func (h *Handler) Payment(c *gin.Context) {
	var req request.PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("invalid login request", "error", err)
		response.ErrorMsg(c, err, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		h.Log.Warn("invalid login request", "error", err)
		response.ErrorMsg(c, err, http.StatusBadRequest)
		return
	}
	// simulate payment process
	time.Sleep(1 * time.Second)

	resData := response.PaymentResponse{
		Status:  "PAID",
		Message: fmt.Sprintf("Success payment for order id %s", req.OrderID),
	}
	response.OK(c, resData)
}
