package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/neatflowcv/seven-skies/internal/pkg/broker"
)

var _ broker.Broker = (*Broker)(nil)

type Broker struct {
	jet              jetstream.JetStream
	stream           jetstream.Stream
	consumerContexts []jetstream.ConsumeContext
}

// NewBroker url example: nats://127.0.0.1:4222
func NewBroker(ctx context.Context, url string, streamName string, topics []string) (*Broker, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("error in connect: %w", err)
	}

	jet, err := jetstream.New(conn)
	if err != nil {
		return nil, fmt.Errorf("error in create jetstream: %w", err)
	}

	// 같은 이름에 설정이 같으면, 이미 존재하는 스트림을 반환한다. error를 반환하지 않음!
	stream, err := jet.CreateStream(ctx, jetstream.StreamConfig{ //nolint:exhaustruct
		Name:     streamName,
		Subjects: topics,
	})
	if err != nil {
		return nil, fmt.Errorf("error in create stream: %w", err)
	}

	return &Broker{
		jet:              jet,
		stream:           stream,
		consumerContexts: nil,
	}, nil
}

func (b *Broker) Close() {
	for _, consumerContext := range b.consumerContexts {
		consumerContext.Stop()
	}
}

func (b *Broker) Publish(ctx context.Context, topic string, message []byte) error {
	_, err := b.jet.Publish(ctx, topic, message)
	if err != nil {
		return fmt.Errorf("error in publish: %w", err)
	}

	return nil
}

func (b *Broker) Subscribe(ctx context.Context, durable string, fn func(topic string, message []byte) error) error {
	consumer, err := b.stream.CreateConsumer(ctx, jetstream.ConsumerConfig{ //nolint:exhaustruct
		Durable:   durable,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("error in create consumer: %w", err)
	}

	consumeContext, err := consumer.Consume(func(msg jetstream.Msg) {
		topic := msg.Subject()
		content := msg.Data()

		err := fn(topic, content)
		if err != nil {
			log.Println("error in callback", err)

			inErr := msg.Nak()
			if inErr != nil {
				log.Println("error in nak", inErr)
			}

			return
		}

		inErr := msg.Ack()
		if inErr != nil {
			log.Println("error in ack", inErr)
		}
	})
	if err != nil {
		return fmt.Errorf("error in subscribe: %w", err)
	}

	b.consumerContexts = append(b.consumerContexts, consumeContext)

	return nil
}
