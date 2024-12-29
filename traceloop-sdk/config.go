package traceloop

import (
	otlp "go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"time"
)

type BackoffConfig struct {
	MaxRetries uint64
}

type Config struct {
	BaseURL         string
	APIKey          string
	TracerName      string
	ServiceName     string
	PollingInterval time.Duration
	BackoffConfig   BackoffConfig
	Exporter        *otlp.Exporter
}
