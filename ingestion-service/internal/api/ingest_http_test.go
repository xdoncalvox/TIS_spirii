package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestHTTPIngestEndpoint(t *testing.T) {
	fakeQueue := &FakeQueue{}
	tracer := trace.NewNoopTracerProvider().Tracer("test")

	handler := &IngestHandler{
		Queue:  fakeQueue,
		Tracer: tracer,
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	body := `{
		"chargerId": "CH-1",
		"timestamp": "2026-01-06T10:00:00Z",
		"metrics": { "voltage": 400 }
	}`

	resp, err := http.Post(
		server.URL+"/events",
		"application/json",
		strings.NewReader(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.StatusCode)
	}
}
