package telemetry

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var tracerProvider *sdktrace.TracerProvider

func InitTracer(serviceName string, endpoint string) sdktrace.Tracer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exporter sdktrace.SpanExporter
	var err error

	if endpoint != "" {
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			log.Fatalf("failed to create OTLP exporter: %v", err)
		}
	} else {
		// Fallback: no-op exporter (important for demos)
		tracerProvider = sdktrace.NewTracerProvider()
		otel.SetTracerProvider(tracerProvider)
		return otel.Tracer(serviceName)
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tracerProvider)

	return otel.Tracer(serviceName)
}

func Shutdown() {
	if tracerProvider == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tracerProvider.Shutdown(ctx); err != nil {
		log.Printf("error shutting down tracer provider: %v", err)
	}
}
