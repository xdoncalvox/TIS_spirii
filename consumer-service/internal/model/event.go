package model

import "time"

type TelemetryEvent struct {
	ChargerID  string                 `json:"chargerId"`
	EventTime  time.Time              `json:"eventTime"`
	IngestTime time.Time              `json:"ingestTime"`
	Metrics    map[string]interface{} `json:"metrics"`
	TraceID    string                 `json:"traceId"`
}
