package usecase

import (
	"context"
	"fmt"

	"github.com/fallinnadim/inventory-worker/internal/domain"
	"github.com/fallinnadim/inventory-worker/internal/port/outbound"
)

type ItemUsecase struct {
	itemRepo  outbound.ItemRepository
	cacheRepo outbound.CacheRepository
}

func NewItemUsecase(itemRepo outbound.ItemRepository, cacheRepo outbound.CacheRepository) *ItemUsecase {
	return &ItemUsecase{itemRepo, cacheRepo}
}

func (i *ItemUsecase) UpdateInventory(ctx context.Context, values []domain.OrderItem) error {
	if err := i.itemRepo.UpdateStocks(ctx, values); err != nil {
		return fmt.Errorf("inventory update failed: %w", err)
	}

	return nil
}

func (i *ItemUsecase) InvalidateCache(ctx context.Context, values []domain.OrderItem) error {
	if err := i.cacheRepo.InvalidateCache(ctx, values); err != nil {
		return fmt.Errorf("cache invalidation failed: %w", err)
	}

	return nil
}
