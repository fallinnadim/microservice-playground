package usecase

import (
	"context"
	"math"
	"time"

	"github.com/fallinnadim/order-service/internal/port/outbound"
)

type RateLimitUsecase struct {
	service    outbound.RateLimitService
	capacity   int
	refillRate float64
}

func NewRateLimitUsecase(service outbound.RateLimitService, capacity int, refillRate float64) *RateLimitUsecase {
	return &RateLimitUsecase{
		service, capacity, refillRate,
	}
}
func (uc *RateLimitUsecase) Allow(ctx context.Context, userID string) (bool, error) {
	key := "ratelimit:" + userID

	tokens, lastRefill, err := uc.service.GetBucket(ctx, key)
	if err != nil {
		return false, err
	}

	now := time.Now().UnixMilli()

	if lastRefill == 0 {
		tokens = float64(uc.capacity)
		lastRefill = now
	}

	elapsed := float64(now-lastRefill) / 1000.0

	tokens = math.Min(
		float64(uc.capacity),
		tokens+(elapsed*uc.refillRate),
	)

	if tokens < 1.0 {
		_ = uc.service.SetBucket(ctx, key, tokens, now)
		return false, nil
	}

	tokens -= 1.0

	err = uc.service.SetBucket(ctx, key, tokens, now)
	if err != nil {
		return false, err
	}

	return true, nil
}
