package usecase

import "github.com/jackc/pgx/v5/pgxpool"

type pingUsecase struct {
	db *pgxpool.Pool
}

func NewPingUsecase(db *pgxpool.Pool) *pingUsecase {
	return &pingUsecase{
		db,
	}
}

func (u *pingUsecase) Ping() string {
	return "pong"
}
