# Telemetry Ingestion System

## What this is
End-to-end telemetry ingestion and observability platform for charger events.
Accepts high-volume time-series events and exposes them via Grafana.

## Architecture
Charger → Ingestion API (OpenTelemetry) → Kafka → TimescaleDB → Grafana

## Components
- ingestion-service: HTTP API emitting OTEL events
- Kafka: durable event buffer
- TimescaleDB: time-series storage
- Grafana: visualization

# ADR
For this task only, we use TimescaleDB to reduce operational complexity. ClickHouse is the intended production store.

## Running locally
docker-compose up

## Event schema
See clickhouse/schema/raw_events.sql