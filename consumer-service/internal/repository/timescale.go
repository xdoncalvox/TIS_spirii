package repository

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"telemetry/consumer-service/internal/model"
)

type Repository struct {
	db *sql.DB
}

func New(dsn string) (*Repository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) InsertEvent(ctx context.Context, e model.TelemetryEvent) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO telemetry_events
		 (charger_id, event_time, ingest_time, metrics, trace_id)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT DO NOTHING`,
		e.ChargerID,
		e.EventTime,
		e.IngestTime,
		e.Metrics,
		e.TraceID,
	)
	return err
}
