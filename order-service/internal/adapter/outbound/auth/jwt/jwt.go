package jwt

import (
	"errors"
	"time"

	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthAdapter struct {
	secret   []byte
	duration string
}

func NewJWTAuthAdapter(secret string, duration string) *JWTAuthAdapter {
	return &JWTAuthAdapter{
		secret:   []byte(secret),
		duration: duration,
	}
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

func (j *JWTAuthAdapter) GenerateToken(userId string) (string, error) {
	duration, _ := time.ParseDuration(j.duration)
	claims := CustomClaims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
