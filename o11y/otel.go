// Package o11y contains functionality that opentelemetry uses to ingest
// telemetry data.
package o11y

import (
	"context"

	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/logger"
	"github.com/thisiserico/golib/oops"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Configure sets up opentelemetry to use golib functionality.
// Specifically, an otel tracer that uses the logger is set.
func Configure(log logger.Log) {
	tracer := tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(newSpanProcessor(log)))
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

func RecordError(span trace.Span, err error) {
	var attributes []attribute.KeyValue
	for _, detail := range oops.Details(err) {
		var attr attribute.KeyValue
		name := detail.Name()

		switch value := detail.Value(); value.(type) {
		case string:
			attr = attribute.String(name, value.(string))

		case int:
			attr = attribute.Int(name, value.(int))

		case bool:
			attr = attribute.Bool(name, value.(bool))
		}

		attributes = append(attributes, attr)
	}

	span.SetAttributes(attributes...)
	span.RecordError(err)
}
