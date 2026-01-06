package queue

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"telemetry/ingestion-service/internal/model"

	"github.com/segmentio/kafka-go"
)

type KafkaQueue struct {
	writer *kafka.Writer
}

func NewKafkaQueue(brokers []string, topic string) *KafkaQueue {
	ensureTopic(brokers, topic)

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

func ensureTopic(brokers []string, topic string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		log.Fatalf("failed to dial kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed to get controller: %v", err)
	}
	controllerConn, err := kafka.DialContext(
		ctx,
		"tcp",
		controller.Host+":"+strconv.Itoa(controller.Port),
	)
	if err != nil {
		log.Fatalf("failed to dial controller: %v", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		// Topic already exists is NOT an error
		log.Printf("topic create result: %v", err)
	} else {
		log.Printf("topic %s created", topic)
	}
}

func (k *KafkaQueue) Close() error {
	return k.writer.Close()
}
