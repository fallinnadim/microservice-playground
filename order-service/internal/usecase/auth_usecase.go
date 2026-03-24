package usecase

import (
	"context"
	"errors"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/port/outbound"
)

type AuthUsecase struct {
	JWTAdapter    outbound.JWTAuthPort
	Argon2Adapter outbound.Argon2Port
	userRepo      outbound.UserRepository
}

func NewAuthUsecase(
	jwtAdapter outbound.JWTAuthPort,
	argon2Adapter outbound.Argon2Port,
	userRepo outbound.UserRepository,
) *AuthUsecase {
	return &AuthUsecase{jwtAdapter, argon2Adapter, userRepo}
}

func (u *AuthUsecase) Login(ctx context.Context, req inbound.LoginInput) (*inbound.LoginOutput, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	ok, err := u.Argon2Adapter.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid password")
	}
	token, err := u.JWTAdapter.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &inbound.LoginOutput{
		Token: token,
	}, nil
}

func (u *AuthUsecase) Register(ctx context.Context, req inbound.LoginInput) error {
	user, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if user != nil {
		return errors.New("email already exist")
	}
	hashedPassword, err := u.Argon2Adapter.GenerateFromPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash")
	}
	errNewUser := u.userRepo.CreateNewUser(ctx, req.Email, hashedPassword)
	if errNewUser != nil {
		return errNewUser
	}
	return nil
}
