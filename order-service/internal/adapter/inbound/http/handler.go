package http

import (
	"log/slog"
	"net/http"

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
	c.JSON(http.StatusOK, gin.H{
		"message": h.pingUC.Ping(),
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req inbound.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid login request", "error", err)
		response.ErrorMsg(c, http.StatusBadRequest, "invalid request")
		return
	}

	res, err := h.authUC.Login(c.Request.Context(), req)
	if err != nil {
		h.log.Warn("login failed", "error", err)
		response.ErrorMsg(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, "Login Successfully", res.Token)
}

func (h *Handler) Register(c *gin.Context) {
	var req inbound.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid register request", "error", err)
		response.ErrorMsg(c, http.StatusBadRequest, "invalid request")
		return
	}

	err := h.authUC.Register(c.Request.Context(), req)
	if err != nil {
		h.log.Warn("register failed", "error", err)
		response.ErrorMsg(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, "Register Successfully", "Register Successfully")
}
