package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/fallinnadim/order-service/internal/port/outbound"
	"github.com/segmentio/kafka-go"
)

type KafkaProducerAdapter struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) outbound.KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		BatchTimeout: 10 * time.Millisecond,
	}
	return &KafkaProducerAdapter{writer: writer}
}

func (k *KafkaProducerAdapter) Publish(ctx context.Context, topic string, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
		Topic: topic,
	}

	if err := k.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write kafka message: %w", err)
	}
	return nil
}
