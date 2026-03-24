package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/request"
	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/response"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	pingUC      inbound.PingUsecase
	authUC      *usecase.AuthUsecase
	rateLimitUC *usecase.RateLimitUsecase
	log         *slog.Logger
}

func NewHandler(
	pingUC inbound.PingUsecase,
	authUC *usecase.AuthUsecase,
	rateLimitUC *usecase.RateLimitUsecase,
	log *slog.Logger,
) *Handler {
	return &Handler{
		pingUC, authUC, rateLimitUC, log,
	}
}

func (h *Handler) Ping(c *gin.Context) {
	h.log.Info("ping called")
	result, err := h.pingUC.Ping(c.Request.Context())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timeout",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid login request", "error", err)
		response.ErrorMsg(c, http.StatusBadRequest, "invalid request")
		return
	}
	input := inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.authUC.Login(c.Request.Context(), input)
	if err != nil {
		h.log.Warn("login failed", "error", err)
		response.ErrorMsg(c, http.StatusUnauthorized, err.Error())
		return
	}
	resData := response.LoginResponse{
		Token: res.Token,
	}
	response.OK(c, "Login Successfully", resData)
}

func (h *Handler) Register(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid register request", "error", err)
		response.ErrorMsg(c, http.StatusBadRequest, "invalid request")
		return
	}

	input := inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.authUC.Register(c.Request.Context(), input)
	if err != nil {
		h.log.Warn("register failed", "error", err)
		response.ErrorMsg(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, "Register Successfully", "Register Successfully")
}
