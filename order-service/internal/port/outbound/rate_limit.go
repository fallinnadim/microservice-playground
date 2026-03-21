package outbound

import "context"

type RateLimitRepository interface {
	GetBucket(ctx context.Context, key string) (tokens int, lastRefill int64, err error)
	SetBucket(ctx context.Context, key string, tokens int, lastRefill int64) error
}
