// Package o11y contains functionality that opentelemetry uses to ingest
// telemetry data.
package o11y

import (
	"context"

	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Configure sets up opentelemetry to use golib functionality.
// Specifically, an otel tracer that uses the logger is set.
func Configure(log logger.Log) {
	tracer := trace.NewTracerProvider(trace.WithSpanProcessor(newSpanProcessor(log)))
	otel.SetTracerProvider(tracer)
}

// Attributes extracts all the known pairs from the context and converts
// them into a slice of attributes.
func Attributes(ctx context.Context) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("correlation_id", kv.CorrelationID(ctx).String()),
		attribute.Bool("is_dry_run", kv.IsDryRun(ctx).Bool()),
	}
}
