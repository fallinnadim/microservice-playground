package ratelimit

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitAdapter struct {
	rdb *redis.Client
}

func NewRateLimitAdapter(rdb *redis.Client) *RateLimitAdapter {
	return &RateLimitAdapter{rdb: rdb}
}

func (r *RateLimitAdapter) GetBucket(ctx context.Context, key string) (float64, int64, error) {
	val, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	if len(val) == 0 {
		return 10, time.Now().Unix(), nil
	}

	tokens, _ := strconv.ParseFloat(val["tokens"], 64)
	lastRefill, _ := strconv.ParseInt(val["last_refill"], 10, 64)

	return tokens, lastRefill, nil
}

func (r *RateLimitAdapter) SetBucket(ctx context.Context, key string, tokens float64, lastRefill int64) error {
	return r.rdb.HSet(ctx, key, map[string]any{
		"tokens":      tokens,
		"last_refill": lastRefill,
	}).Err()
}
