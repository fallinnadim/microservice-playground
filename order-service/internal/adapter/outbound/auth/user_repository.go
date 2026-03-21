package auth

import (
	"context"

	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/fallinnadim/order-service/internal/port/outbound"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) outbound.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateNewUser(ctx context.Context, email string, password string) error {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT user_id, email
		FROM users
		WHERE email = $1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
