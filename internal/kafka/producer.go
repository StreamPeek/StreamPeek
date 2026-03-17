package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// ProduceResult contains the result of a successful produce operation.
type ProduceResult struct {
	Topic      string
	Partition  int32
	Offset     int64
	BrokerID   int32
	DurationMs int64
}

// Client wraps a franz-go kgo.Client.
type Client struct {
	client *kgo.Client
}

// NewClient initializes a new Kafka client connected to the given seed brokers.
func NewClient(brokers string) (*Client, error) {
	b := strings.Split(brokers, ",")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(b...),
		// Default to at-least-once semantics for now
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// Close closes the underlying Kafka client.
func (c *Client) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// ProduceSync synchronously produces a record to Kafka.
func (c *Client) ProduceSync(ctx context.Context, topic string, key, value []byte, headers []kgo.RecordHeader) (*ProduceResult, error) {
	record := &kgo.Record{
		Topic:   topic,
		Key:     key,
		Value:   value,
		Headers: headers,
	}

	start := time.Now()

	// Synchronously produce
	res := c.client.ProduceSync(ctx, record)

	duration := time.Since(start).Milliseconds()

	err := res.FirstErr()
	if err != nil {
		return nil, fmt.Errorf("failed to produce record: %w", err)
	}

	// Assuming success, extract the details from the first (and only) result.
	if len(res) == 0 {
		return nil, fmt.Errorf("produce succeeded but no results returned")
	}

	pr := res[0]
	return &ProduceResult{
		Topic:      pr.Record.Topic,
		Partition:  pr.Record.Partition,
		Offset:     pr.Record.Offset,
		BrokerID:   1, // Stub for now or attempt extraction if metadata available.
		DurationMs: duration,
	}, nil
}
