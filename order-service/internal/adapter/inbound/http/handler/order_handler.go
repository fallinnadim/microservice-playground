package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/request"
	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/response"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	validate *validator.Validate
	orderUC  inbound.OrderUsecase
	log      *slog.Logger
}

func NewOrderHandler(
	validate *validator.Validate,
	orderUC inbound.OrderUsecase,
	log *slog.Logger,
) *OrderHandler {
	return &OrderHandler{
		validate, orderUC, log,
	}
}

func (h *OrderHandler) Order(c *gin.Context) {
	var req request.OrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid order request", "error", err)
		response.ErrorMsg(c, err, http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(&req); err != nil {
		h.log.Warn("invalid order request", "error", err)
		response.ErrorMsg(c, err, http.StatusBadRequest)
		return
	}
	var items []inbound.ItemInput
	for _, v := range req.Items {
		item := inbound.ItemInput{
			Id:      v.Id,
			Ammount: v.Ammount,
		}
		items = append(items, item)
	}
	input := inbound.OrderInput{
		UserId: req.UserId,
		Items:  items,
	}
	result, err := h.orderUC.Order(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientStock) {
			response.ErrorMsg(c, err, http.StatusConflict)
			return
		}
		if errors.Is(err, usecase.ErrFailedPayment) {
			response.ErrorMsg(c, err, http.StatusPaymentRequired)
			return
		}
		response.ErrorMsg(c, err, http.StatusInternalServerError)
		return
	}
	response.Created(c, result)
}
