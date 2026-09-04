// Package o11y contains functionality that opentelemetry uses to ingest
// telemetry data.
package o11y

import (
	"context"

	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger"
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
	attrs := make([]attribute.KeyValue, 0, 5)
	for _, attr := range kv.AllAttributes(ctx) {
		attrs = append(attrs, attribute.Any(attr.Name(), attr.Value()))
	}

	return attrs
}
