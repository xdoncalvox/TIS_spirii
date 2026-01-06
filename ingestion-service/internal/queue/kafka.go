package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"telemetry/internal/model"
)

type KafkaQueue struct {
	writer *kafka.Writer
}

func NewKafkaQueue(brokers []string, topic string) *KafkaQueue {
	return &KafkaQueue{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{}, // key-based partitioning
			RequiredAcks: kafka.RequireOne,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (k *KafkaQueue) Publish(ctx context.Context, key string, event model.TelemetryEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: bytes,
		Time:  event.EventTime,
	}

	return k.writer.WriteMessages(ctx, msg)
}

func (k *KafkaQueue) Close() error {
	return k.writer.Close()
}
