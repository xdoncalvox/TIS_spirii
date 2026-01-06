package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"telemetry/consumer-service/internal/model"
)

type Consumer struct {
	Reader  MessageReader
	Tracer  trace.Tracer
	Process func(ctx context.Context, event model.TelemetryEvent)
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("kafka read error: %v", err)
			continue
		}

		c.handleMessage(ctx, msg)
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg kafka.Message) {
	ctx, span := c.Tracer.Start(
		ctx,
		"kafka.consume",
		trace.WithAttributes(
			attribute.String("charger.id", string(msg.Key)),
		),
	)
	defer span.End()

	var event model.TelemetryEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		span.RecordError(err)
		return
	}

	c.Process(ctx, event)
}
