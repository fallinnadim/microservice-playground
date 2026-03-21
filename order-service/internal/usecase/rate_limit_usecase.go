package usecase

import (
	"context"
	"time"

	"github.com/fallinnadim/order-service/internal/port/outbound"
)

type RateLimitUsecase struct {
	repo       outbound.RateLimitRepository
	capacity   int
	refillRate int
}

func NewRateLimitUsecase(repo outbound.RateLimitRepository, capacity, refillRate int) *RateLimitUsecase {
	return &RateLimitUsecase{
		repo:       repo,
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (uc *RateLimitUsecase) Allow(ctx context.Context, userID string) (bool, error) {
	key := "ratelimit:" + userID

	tokens, lastRefill, err := uc.repo.GetBucket(ctx, key)
	if err != nil {
		return false, err
	}

	now := time.Now().Unix()

	elapsed := now - lastRefill
	refilled := int(elapsed) * uc.refillRate

	if refilled > 0 {
		tokens = min(uc.capacity, tokens+refilled)
		lastRefill = now
	}

	if tokens <= 0 {
		return false, nil
	}

	tokens--

	err = uc.repo.SetBucket(ctx, key, tokens, lastRefill)
	if err != nil {
		return false, err
	}

	return true, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
