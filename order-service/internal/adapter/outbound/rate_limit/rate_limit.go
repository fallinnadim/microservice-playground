package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRedisRepository struct {
	rdb *redis.Client
}

func NewRateLimitRedisRepository(rdb *redis.Client) *RateLimitRedisRepository {
	return &RateLimitRedisRepository{rdb: rdb}
}

func (r *RateLimitRedisRepository) GetBucket(ctx context.Context, key string) (int, int64, error) {
	val, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	if len(val) == 0 {
		return 10, time.Now().Unix(), nil
	}

	tokens, _ := strconv.Atoi(val["tokens"])
	lastRefill, _ := strconv.ParseInt(val["last_refill"], 10, 64)

	return tokens, lastRefill, nil
}

func (r *RateLimitRedisRepository) SetBucket(ctx context.Context, key string, tokens int, lastRefill int64) error {
	return r.rdb.HSet(ctx, key, map[string]interface{}{
		"tokens":      tokens,
		"last_refill": lastRefill,
	}).Err()
}
