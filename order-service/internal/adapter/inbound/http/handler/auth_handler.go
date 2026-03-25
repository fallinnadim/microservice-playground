package handler

import (
	"log/slog"
	"net/http"

	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/request"
	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/response"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	validate *validator.Validate
	AuthUC   *usecase.AuthUsecase
	Log      *slog.Logger
}

func NewAuthHandler(
	validate *validator.Validate,
	authUC *usecase.AuthUsecase,
	log *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		validate, authUC, log,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

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
	input := inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := h.AuthUC.Login(c.Request.Context(), input)
	if err != nil {
		h.Log.Warn("login failed", "error", err)
		response.ErrorMsg(c, err, http.StatusUnauthorized)
		return
	}
	resData := response.LoginResponse{
		Token: res.Token,
	}
	response.OK(c, resData)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("invalid register request", "error", err)
		response.ErrorMsg(c, err, http.StatusBadRequest)
		return
	}

	input := inbound.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.AuthUC.Register(c.Request.Context(), input)
	if err != nil {
		h.Log.Warn("register failed", "error", err)
		response.ErrorMsg(c, err, http.StatusUnauthorized)
		return
	}
	response.OK(c, "Register Successfully")
}
