package item

import (
	"context"
	"fmt"

	"github.com/fallinnadim/inventory-worker/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type itemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *itemRepository {
	return &itemRepository{db: db}
}

func (i *itemRepository) UpdateStocks(ctx context.Context, items []domain.OrderItem) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		query := `UPDATE items SET stock = stock - $1 WHERE id = $2 AND stock >= $1`

		cmd, err := tx.Exec(ctx, query, item.Quantity, item.ItemID)
		if err != nil {
			return fmt.Errorf("failed to update item %s: %w", item.ItemID, err)
		}

		if cmd.RowsAffected() == 0 {
			return fmt.Errorf("insufficient stock or item not found: %s", item.ItemID)
		}
	}

	return tx.Commit(ctx)
}
