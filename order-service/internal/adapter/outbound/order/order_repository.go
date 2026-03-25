package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *orderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) CreateOrder(ctx context.Context, orderInput OrderInput) (string, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var orderID string

	queryOrder := `
		INSERT INTO orders (user_id, status)
		VALUES ($1, $2)
		RETURNING id
	`

	err = tx.QueryRow(ctx, queryOrder, orderInput.UserId, orderInput.Status).Scan(&orderID)
	if err != nil {
		return "", err
	}

	queryItem := `
		INSERT INTO order_items (order_id, item_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range orderInput.Items {
		_, err := tx.Exec(ctx, queryItem,
			orderID,
			item.ItemID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return orderID, nil
}
func (o *orderRepository) UpdateOrder(ctx context.Context, orderId, newStatus string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`

	res, err := o.db.Exec(ctx, query, newStatus, orderId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
