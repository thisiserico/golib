// Package trace let clients handle tracing segments.
package trace

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/thisiserico/golib/constant"
)

// Segment encapsulates a tracing span.
type Segment struct {
	name      string
	createdAt time.Time
	span      opentracing.Span
}

// NewSegment initializes a new tracing segment.
func NewSegment(ctx context.Context, name string) *Segment {
	createdAt := time.Now()

	var parentCtx opentracing.SpanContext
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	return &Segment{
		name:      name,
		createdAt: createdAt,
		span: opentracing.StartSpan(
			name,
			opentracing.StartTime(createdAt),
			opentracing.ChildOf(parentCtx),
			opentracing.Tags{
				"name": name,
			},
		),
	}
}

// Log records different events.
func (s *Segment) Log(key constant.Key, value constant.Value) {
	s.span.LogFields(log.Object(string(key), value))
}

// Finish finalizes the segment span.
func (s *Segment) Finish(err *error) {
	if err != nil && *err != nil {
		sErr := *err
		s.span.SetTag("error", true)
		s.span.SetTag("error_msg", sErr.Error())
	}

	s.span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: time.Now(),
	})
}
