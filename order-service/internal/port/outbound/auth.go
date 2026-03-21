package outbound

import "github.com/fallinnadim/order-service/internal/domain"

type AuthService interface {
	ValidateToken(token string) (*domain.Claims, error)
}
