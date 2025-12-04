//go:build integration

package nats_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/neatflowcv/seven-skies/internal/pkg/broker/nats"
	"github.com/stretchr/testify/require"
)

func TestNewBroker(t *testing.T) {
	broker, err := nats.NewBroker(
		t.Context(),
		"nats://127.0.0.1:4222",
		"SEVEN_SKIES_TEST_STREAM",
		[]string{
			"SEVEN_SKIES_TEST_SUBJECT.>",
		},
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		broker.Close()
	})

	err = broker.Subscribe(t.Context(), "SEVEN_SKIES_TEST_DURABLE", func(topic string, message []byte) error {
		fmt.Println(topic, string(message))
		return nil
	})
	require.NoError(t, err)

	err = broker.Publish(t.Context(), "SEVEN_SKIES_TEST_SUBJECT.data", []byte("hello"))
	require.NoError(t, err)

	time.Sleep(2 * time.Second)
}
