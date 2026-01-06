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
Kafka will run with a single broker for this task, in a production setup more brokers (ideally 3 as minimum) would be required.

Database credentials are shared in plain text for simplicity of this challenge and by no means reflect the correct/suggested behavior, proper secret management should be in place in a real scenario.

## Running locally
docker-compose up

## Event schema
See clickhouse/schema/raw_events.sql