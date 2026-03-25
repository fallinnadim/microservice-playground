package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/fallinnadim/order-service/internal/adapter/outbound/item"
	"github.com/fallinnadim/order-service/internal/adapter/outbound/order"
	"github.com/fallinnadim/order-service/internal/domain"
	"github.com/fallinnadim/order-service/internal/port/inbound"
	"github.com/fallinnadim/order-service/internal/port/outbound"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock for")
	ErrFailedPayment     = errors.New("payment failure for")
)

type orderUsecase struct {
	itemRepo      outbound.ItemRepository
	cacheRepo     outbound.ItemCacheRepository
	orderRepo     outbound.OrderRepository
	paymentClient outbound.PaymentService
	kafkaProducer outbound.KafkaProducer
}

func NewOrderUsecase(
	itemRepo outbound.ItemRepository,
	cacheRepo outbound.ItemCacheRepository,
	orderRepo outbound.OrderRepository,
	paymentClient outbound.PaymentService,
	kafkaProducer outbound.KafkaProducer,

) *orderUsecase {
	return &orderUsecase{
		itemRepo, cacheRepo, orderRepo, paymentClient, kafkaProducer,
	}
}

func (o *orderUsecase) Order(ctx context.Context, orderInput inbound.OrderInput) (string, error) {
	var itemIds []string
	for _, v := range orderInput.Items {
		itemIds = append(itemIds, v.Id)
	}
	keys := buildKeys(itemIds)

	cachedMap, _ := o.cacheRepo.GetItems(ctx, keys)

	var missingIDs []string
	for _, id := range itemIds {
		key := buildKey(id)
		if _, ok := cachedMap[key]; !ok {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		dbItems, err := o.itemRepo.FindByIds(ctx, missingIDs)
		if err != nil {
			return "", err
		}

		toCache := make(map[string]item.ItemCache)
		for _, itemz := range dbItems {
			key := buildKey(itemz.ID)

			toCache[key] = item.ItemCache{
				ID:    itemz.ID,
				Name:  itemz.Name,
				Price: itemz.Price,
				Stock: itemz.Stock,
			}
		}
		_ = o.cacheRepo.SetItems(ctx, toCache, 5*time.Minute)

		maps.Copy(cachedMap, toCache)
	}

	var items []domain.Item
	for _, id := range itemIds {
		key := buildKey(id)

		if cachedItem, ok := cachedMap[key]; ok {
			items = append(items, domain.Item{
				ID:    cachedItem.ID,
				Name:  cachedItem.Name,
				Price: cachedItem.Price,
				Stock: cachedItem.Stock,
			})
		}
	}
	inputMap := make(map[string]int)

	for _, input := range orderInput.Items {
		inputMap[input.Id] = input.Ammount
	}
	for _, item := range items {
		requested := inputMap[item.ID]

		if requested > item.Stock {
			return "", fmt.Errorf("%w: item_id=%s", ErrInsufficientStock, item.ID)
		}
	}
	var itemz []domain.OrderItem
	var totalAmount int
	for _, item := range items {
		requested := inputMap[item.ID]
		newItemz := domain.OrderItem{
			ItemID:   item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: requested,
		}
		itemz = append(itemz, newItemz)
		totalAmount += (newItemz.Price * newItemz.Quantity)
	}
	input := order.OrderInput{
		UserId: orderInput.UserId,
		Status: "PENDING",
		Items:  itemz,
	}
	orderId, _ := o.orderRepo.CreateOrder(ctx, input)
	fmt.Println(totalAmount)
	paymentResp, err := o.paymentClient.Pay(ctx, order.PaymentRequest{
		OrderID: orderId,
		UserID:  orderInput.UserId,
		Amount:  totalAmount,
	})

	if err != nil {
		return "", err
	}
	if paymentResp.Status == "PAID" {
		err = o.orderRepo.UpdateOrder(ctx, orderId, "PAID")
		if err != nil {
			return "", err
		}
		value, _ := json.Marshal(itemz)
		errz := o.kafkaProducer.Publish(ctx, "order.created", []byte(orderId), value)
		if errz != nil {
			fmt.Printf("error kafka topic %v", errz)
		}
		return fmt.Sprintf("Success payment for order id : %s", orderId), nil

	} else {
		err = o.orderRepo.UpdateOrder(ctx, orderId, "FAILED")
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%w: order id : %s", ErrFailedPayment, orderId)
	}
}

func buildKey(id string) string {
	return "item:" + id
}

func buildKeys(ids []string) []string {
	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = buildKey(id)
	}
	return keys
}
