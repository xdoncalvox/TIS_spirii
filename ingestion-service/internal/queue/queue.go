package queue

import "context"
import "telemetry/internal/model"

type EventQueue interface {
	Publish(ctx context.Context, key string, event model.TelemetryEvent) error
	Close() error
}
