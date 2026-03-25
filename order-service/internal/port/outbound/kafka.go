package outbound

import "context"

type KafkaProducer interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}
