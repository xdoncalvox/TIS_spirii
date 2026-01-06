package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"

	"telemetry/ingestion-service/internal/model"
	"telemetry/ingestion-service/internal/queue"
)

type IngestHandler struct {
	Queue  queue.EventQueue
	Tracer trace.Tracer
}

func (h *IngestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ChargerID string                 `json:"chargerId"`
		Timestamp time.Time              `json:"timestamp"`
		Metrics   map[string]interface{} `json:"metrics"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if payload.ChargerID == "" || payload.Timestamp.IsZero() {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	ctx, span := h.Tracer.Start(r.Context(), "ingest_event")
	defer span.End()

	event := model.TelemetryEvent{
		ChargerID:  payload.ChargerID,
		EventTime:  payload.Timestamp,
		IngestTime: time.Now().UTC(),
		Metrics:    payload.Metrics,
		TraceID:    span.SpanContext().TraceID().String(),
	}

	if err := h.Queue.Publish(ctx, payload.ChargerID, event); err != nil {
		http.Error(w, "queue failure", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
