package jwt

import (
	"errors"

	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthAdapter struct {
	secret []byte
}

func NewJWTAuthAdapter(secret string) *JWTAuthAdapter {
	return &JWTAuthAdapter{secret: []byte(secret)}
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *JWTAuthAdapter) ValidateToken(tokenStr string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &domain.Claims{
		UserID: claims.UserID,
	}, nil
}
