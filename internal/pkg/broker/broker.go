package broker

import "context"

type Broker interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, durable string, fn func(topic string, message []byte) error) error
}
