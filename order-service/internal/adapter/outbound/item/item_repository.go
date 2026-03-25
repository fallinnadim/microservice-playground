package item

import (
	"context"

	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type itemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *itemRepository {
	return &itemRepository{db: db}
}

func (i *itemRepository) FindByIds(ctx context.Context, ids []string) ([]*domain.Item, error) {
	query := `
		SELECT id, name, price, stock
		FROM items
		WHERE id = ANY($1)
	`
	rows, err := i.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Stock); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
