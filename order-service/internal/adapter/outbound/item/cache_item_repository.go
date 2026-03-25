package item

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type itemCacheRepository struct {
	rdb *redis.Client
}

func NewItemCacheRepository(rdb *redis.Client) *itemCacheRepository {
	return &itemCacheRepository{rdb}
}

func (i *itemCacheRepository) GetItems(ctx context.Context, keys []string) (map[string]ItemCache, error) {
	result := make(map[string]ItemCache)

	values, err := i.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	for idx, val := range values {
		if val == nil {
			continue
		}

		strVal, ok := val.(string)
		if !ok {
			continue
		}

		var item ItemCache
		if err := json.Unmarshal([]byte(strVal), &item); err != nil {
			continue
		}

		result[keys[idx]] = item
	}

	return result, nil
}
func (i *itemCacheRepository) SetItems(ctx context.Context, items map[string]ItemCache, ttl time.Duration) error {
	pipe := i.rdb.Pipeline()

	for key, item := range items {
		bytes, err := json.Marshal(item)
		if err != nil {
			continue
		}

		pipe.Set(ctx, key, bytes, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}
