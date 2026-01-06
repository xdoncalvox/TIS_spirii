CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE IF NOT EXISTS telemetry_events (
  charger_id   TEXT        NOT NULL,
  event_time   TIMESTAMPTZ NOT NULL,
  ingest_time  TIMESTAMPTZ NOT NULL,
  metrics      JSONB       NOT NULL,
  trace_id     TEXT,

  -- idempotency key
  PRIMARY KEY (charger_id, event_time)
);

SELECT create_hypertable(
  'telemetry_events',
  'event_time',
  if_not_exists => TRUE
);
