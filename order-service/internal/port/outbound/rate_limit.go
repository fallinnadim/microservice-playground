package outbound

import "context"

type RateLimitService interface {
	GetBucket(ctx context.Context, key string) (tokens float64, lastRefill int64, err error)
	SetBucket(ctx context.Context, key string, tokens float64, lastRefill int64) error
}
