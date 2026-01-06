package api

import (
	"context"
	"telemetry/internal/model"
)
type FakeQueue struct {
	Published []model.TelemetryEvent
}

func (f *FakeQueue) Publish(ctx context.Context, key string, event model.TelemetryEvent) error {
	f.Published = append(f.Published, event)
	return nil
}

func (f *FakeQueue) Close() error {
	return nil
}

func TestIngestHandler_AcceptsValidEvent(t *testing.T) {
	fakeQueue := &FakeQueue{}
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	handler := &IngestHandler{
		Queue:  fakeQueue,
		Tracer: tracer,
	}

	body := `{
		"chargerId": "CH-1",
		"timestamp": "2026-01-06T10:00:00Z",
		"metrics": { "voltage": 400 }
	}`

	req := httptest.NewRequest("POST", "/events", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}

	if len(fakeQueue.Published) != 1 {
		t.Fatalf("expected 1 event published")
	}

	event := fakeQueue.Published[0]

	if event.ChargerID != "CH-1" {
		t.Fatalf("chargerId mismatch")
	}

	if event.IngestTime.IsZero() {
		t.Fatalf("ingestTime not set")
	}
}
