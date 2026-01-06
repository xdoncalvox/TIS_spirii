package main

import (
	"log"
	"net/http"

	"telemetry/ingestion-service/internal/api"
	"telemetry/ingestion-service/internal/config"
	"telemetry/ingestion-service/internal/queue"
	"telemetry/ingestion-service/internal/telemetry"
)

func main() {
	cfg := config.Load()

	tracer := telemetry.InitTracer(cfg.ServiceName, cfg.OtelEndpoint)
	defer telemetry.Shutdown()

	kafkaQueue := queue.NewKafkaQueue(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer kafkaQueue.Close()

	handler := &api.IngestHandler{
		Queue:  kafkaQueue,
		Tracer: tracer,
	}

	http.Handle("/events", handler)

	log.Println("ingestion service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
