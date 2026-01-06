package consumer

import (
	"context"
	"errors"
	"telemetry/consumer-service/internal/model"
	"testing"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type FakeReader struct {
	Messages []kafka.Message
	index    int
}

func (f *FakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if f.index >= len(f.Messages) {
		return kafka.Message{}, errors.New("no more messages")
	}
	msg := f.Messages[f.index]
	f.index++
	return msg, nil
}

func (f *FakeReader) Close() error {
	return nil
}

func TestConsumerProcessesValidMessage(t *testing.T) {
	fakeReader := &FakeReader{
		Messages: []kafka.Message{
			{
				Key: []byte("CH-1"),
				Value: []byte(`{
					"chargerId": "CH-1",
					"eventTime": "2026-01-06T12:00:00Z",
					"metrics": { "voltage": 400 }
				}`),
			},
		},
	}

	called := false

	consumer := &Consumer{
		Reader: fakeReader,
		Tracer: trace.NewNoopTracerProvider().Tracer("test"),
		Process: func(ctx context.Context, event model.TelemetryEvent) {
			called = true
			if event.ChargerID != "CH-1" {
				t.Fatalf("wrong chargerId")
			}
		},
	}

	// Run only ONE iteration
	msg, _ := fakeReader.ReadMessage(context.Background())
	consumer.handleMessage(context.Background(), msg)

	if !called {
		t.Fatalf("expected process to be called")
	}
}

func TestConsumerSkipsInvalidJSON(t *testing.T) {
	fakeReader := &FakeReader{
		Messages: []kafka.Message{
			{Value: []byte(`not-json`)},
		},
	}

	called := false

	consumer := &Consumer{
		Reader: fakeReader,
		Tracer: trace.NewNoopTracerProvider().Tracer("test"),
		Process: func(ctx context.Context, event model.TelemetryEvent) {
			called = true
		},
	}

	msg, _ := fakeReader.ReadMessage(context.Background())
	consumer.handleMessage(context.Background(), msg)

	if called {
		t.Fatalf("process should not be called on invalid JSON")
	}
}
