package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"

	"telemetry/consumer-service/internal/consumer"
	"telemetry/consumer-service/internal/telemetry"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	if topic == "" || groupID == "" {
		log.Fatal("Kafka configuration missing")
	}

	tracer := telemetry.InitTracer("consumer-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	defer telemetry.Shutdown()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	defer reader.Close()

	c := &consumer.Consumer{
		Reader: reader,
		Tracer: tracer,
	}

	c.Process = c.ProcessEvent

	log.Println("consumer-service started")
	c.Run(context.Background())
}
