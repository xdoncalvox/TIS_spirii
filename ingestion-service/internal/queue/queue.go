package queue

import (
	"context"
	"telemetry/ingestion-service/internal/model"
)

type EventQueue interface {
	Publish(ctx context.Context, key string, event model.TelemetryEvent) error
	Close() error
}
