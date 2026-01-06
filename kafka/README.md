# Kafka

This folder documents how Kafka is used within the telemetry platform.

Kafka is the central event backbone that decouples ingestion from downstream
processing (consumers, storage, aggregation).

Kafka is started and managed **only** via the root `docker-compose.yml`.
There is no standalone Kafka stack in this repository.

---

## Purpose

Kafka is used to:

- Buffer high-volume telemetry events
- Preserve ordering per `chargerId`
- Decouple producers (ingestion-service) from consumers
- Allow multiple independent consumers in the future

Kafka is **not** used for:
- Storage
- Analytics
- Querying
- Business logic

---

## Mode of operation

- Single broker (demo-only)
- KRaft mode (no ZooKeeper)
- Auto topic creation enabled
- No authentication (demo-only)

This configuration is intentionally minimal and **not production-ready**.

---

## Services that depend on Kafka

- ingestion-service (producer)
- consumer-service (consumer)

Kafka must be running before either service starts.

---

## Data flow

Charger → ingestion-service → Kafka → consumer-service → (future storage)

---

## Notes

- Topic definitions are documented in `topics.md`
- Partitioning strategy is intentional and documented
- For production, this setup must be replaced with a multi-broker cluster
