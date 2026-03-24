package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pingUsecase struct {
	db *pgxpool.Pool
}

func NewPingUsecase(db *pgxpool.Pool) *pingUsecase {
	return &pingUsecase{
		db,
	}
}

func (u *pingUsecase) Ping(ctx context.Context) (string, error) {
	select {
	case <-time.After(10 * time.Second):
		return "pong", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
