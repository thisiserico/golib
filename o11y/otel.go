// Package o11y contains functionality that opentelemetry uses to ingest
// telemetry data.
package o11y

import (
	"github.com/thisiserico/golib/v2/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Configure sets up opentelemetry to use golib functionality.
// Specifically, an otel tracer that uses the logger is set.
func Configure(log logger.Log) {
	tracer := trace.NewTracerProvider(trace.WithSpanProcessor(newSpanProcessor(log)))
	otel.SetTracerProvider(tracer)
}
