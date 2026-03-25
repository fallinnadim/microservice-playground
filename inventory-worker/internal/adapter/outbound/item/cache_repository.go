package item

import (
	"context"
	"fmt"

	"github.com/fallinnadim/inventory-worker/internal/domain"
	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) *cacheRepository {
	return &cacheRepository{rdb}
}

func (i *cacheRepository) InvalidateCache(ctx context.Context, items []domain.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	itemIds := make([]string, len(items))
	for idx, v := range items {
		itemIds[idx] = v.ItemID
	}

	keys := buildKeys(itemIds)
	err := i.rdb.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate redis cache: %w", err)
	}

	return nil
}

func buildKeys(ids []string) []string {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = buildKey(id)
	}
	return keys
}

func buildKey(id string) string {
	return "item:" + id
}
