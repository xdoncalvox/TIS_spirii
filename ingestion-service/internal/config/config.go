package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	ServiceName string
	HTTPPort    string

	KafkaBrokers []string
	KafkaTopic   string

	OtelEndpoint string
}

func Load() Config {
	cfg := Config{
		ServiceName: getEnv("SERVICE_NAME", "ingestion-service"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),

		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", ""), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", ""),

		OtelEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
	}

	if len(cfg.KafkaBrokers) == 0 || cfg.KafkaBrokers[0] == "" {
		log.Fatal("KAFKA_BROKERS is required")
	}

	if cfg.KafkaTopic == "" {
		log.Fatal("KAFKA_TOPIC is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
