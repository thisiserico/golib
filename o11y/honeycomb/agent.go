package honeycomb

import (
	"context"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/trace"
	"github.com/thisiserico/golib/v2/cntxt"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/o11y"
)

const (
	defaultSpanName = "new"
	existingSpan    = contextKey("existing_span")
)

type contextKey string

type agent struct{}

// Agent prepares the dummy agent so it can be registered.
func Agent(dataset, writeKey string) o11y.Agent {
	beeline.Init(beeline.Config{
		Dataset:  dataset,
		WriteKey: writeKey,
	})

	return &agent{}
}

func (a *agent) StartSpan(ctx context.Context, name string) (context.Context, o11y.Span) {
	ctx, span := beeline.StartSpan(ctx, name)

	return ctx, newSpan(ctx, span)
}

func (a *agent) GetSpan(ctx context.Context) o11y.Span {
	span := trace.GetSpanFromContext(ctx)
	return newSpan(ctx, span)
}

func (a *agent) Flush() {
	beeline.Close()
}

type span struct {
	traceSpan *trace.Span
}

func newSpan(ctx context.Context, s *trace.Span) o11y.Span {
	span := &span{
		traceSpan: s,
	}

	var p kv.Pair
	p = cntxt.BuildID(ctx)
	if value := p.String(); value != "" {
		span.AddPair(ctx, p)
	}
	p = cntxt.ServiceHost(ctx)
	if value := p.String(); value != "" {
		span.AddPair(ctx, p)
	}
	p = cntxt.ServiceName(ctx)
	if value := p.String(); value != "" {
		span.AddPair(ctx, p)
	}
	p = cntxt.CorrelationID(ctx)
	if value := p.String(); value != "" {
		span.AddPair(ctx, p)
	}
	p = cntxt.TriggeredBy(ctx)
	if value := p.String(); value != "" {
		span.AddPair(ctx, p)
	}
	p = cntxt.IsDryRun(ctx)
	span.AddPair(ctx, p)

	return span
}

func (s *span) AddPair(ctx context.Context, pair kv.Pair) {
	s.traceSpan.AddField(pair.Name(), pair.Value())
}

func (s *span) Complete() {
	s.traceSpan.Send()
}
