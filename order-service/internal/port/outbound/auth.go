package outbound

import (
	"context"

	"github.com/fallinnadim/order-service/internal/domain"
)

type AuthTokenService interface {
	ValidateToken(token string) (*domain.Claims, error)
	GenerateToken(userID string) (string, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateNewUser(ctx context.Context, email string, password string) error
}
