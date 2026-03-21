package usecase

import (
	"context"
	"errors"

	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/port/outbound"
)

type AuthUsecase struct {
	authService outbound.AuthTokenService
	userRepo    outbound.UserRepository
}

func NewAuthUsecase(
	authService outbound.AuthTokenService,
	userRepo outbound.UserRepository,
) *AuthUsecase {
	return &AuthUsecase{authService, userRepo}
}

func (u *AuthUsecase) ValidateToken(token string) (*domain.Claims, error) {
	return u.authService.ValidateToken(token)
}

func (u *AuthUsecase) GenerateToken(userId string) (string, error) {
	return u.authService.GenerateToken(userId)
}

func (u *AuthUsecase) Login(ctx context.Context, req inbound.LoginRequest) (*inbound.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	token, err := u.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &inbound.LoginResponse{
		Token: token,
	}, nil
}

func (u *AuthUsecase) Register(ctx context.Context, req inbound.RegisterRequest) error {
	user, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if user != nil {
		return errors.New("email already exist")
	}
	errNewUser := u.userRepo.CreateNewUser(ctx, req.Email, req.Password)
	if errNewUser != nil {
		return errNewUser
	}
	return nil
}
