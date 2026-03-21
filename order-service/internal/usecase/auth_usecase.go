package usecase

import (
	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/fallinnadim/order-service/internal/port/outbound"
)

type AuthUsecase struct {
	authService outbound.AuthService
}

func NewAuthUsecase(authService outbound.AuthService) *AuthUsecase {
	return &AuthUsecase{authService: authService}
}

func (u *AuthUsecase) ValidateToken(token string) (*domain.Claims, error) {
	return u.authService.ValidateToken(token)
}
