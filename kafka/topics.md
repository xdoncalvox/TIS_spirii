# Kafka Topics

This document defines the Kafka topics used by the telemetry platform.

These definitions are **logical contracts**. Topics are auto-created at runtime
by Kafka for demo purposes.

---

## Topic: `charger.telemetry.raw`

### Description
Raw telemetry events emitted by chargers and ingested by the ingestion-service.

This topic contains **unprocessed, append-only events**.

---

### Producer
- ingestion-service

---

### Consumers
- consumer-service (logging / persistence)
- future: aggregation services
- future: monitoring / alerting pipelines

---

### Message Key
- chargerId


**Why:**
- Guarantees ordering per charger
- Enables partition-affinity processing
- Prevents cross-charger interleaving issues

Messages without a key are considered invalid usage.

---

### Message Value
JSON-encoded telemetry event.

Example:

```json
{
  "chargerId": "CH-001",
  "eventTime": "2026-01-06T12:00:00Z",
  "ingestTime": "2026-01-06T12:00:01Z",
  "metrics": {
    "voltage": 401.2,
    "current": 31.8,
    "temperature": 46.1,
    "status": "charging"
  },
  "traceId": "abc123"
}
