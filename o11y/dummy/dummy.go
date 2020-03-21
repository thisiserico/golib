package dummy

import (
	"context"
	"fmt"

	"github.com/thisiserico/golib/v2/cntxt"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger"
	"github.com/thisiserico/golib/v2/o11y"
)

const defaultSpanName = "new"

var existingSpan = contextKey("existing_span")

type contextKey string

type agent struct {
	log logger.Log
}

// Agent prepares the dummy agent so it can be registered.
func Agent(log logger.Log) o11y.Agent {
	return &agent{
		log: log,
	}
}

func (a *agent) StartSpan(ctx context.Context, name string) (context.Context, o11y.Span) {
	return newSpan(ctx, a.log, name)
}

func (a *agent) GetSpan(ctx context.Context) o11y.Span {
	if val := ctx.Value(existingSpan); val != nil {
		return val.(*span)
	}

	_, s := newSpan(ctx, a.log, defaultSpanName)
	return s
}

func (a *agent) Flush() {}

type span struct {
	log logger.Log

	givenName string
	fields    []kv.Pair
}

func newSpan(ctx context.Context, log logger.Log, name string) (context.Context, o11y.Span) {
	fields := make([]kv.Pair, 0, 6)

	var p kv.Pair
	p = cntxt.BuildID(ctx)
	if value := p.String(); value != "" {
		fields = append(fields, p)
	}
	p = cntxt.ServiceHost(ctx)
	if value := p.String(); value != "" {
		fields = append(fields, p)
	}
	p = cntxt.ServiceName(ctx)
	if value := p.String(); value != "" {
		fields = append(fields, p)
	}
	p = cntxt.CorrelationID(ctx)
	if value := p.String(); value != "" {
		fields = append(fields, p)
	}
	p = cntxt.TriggeredBy(ctx)
	if value := p.String(); value != "" {
		fields = append(fields, p)
	}
	p = cntxt.IsDryRun(ctx)
	fields = append(fields, p)

	s := &span{
		log:       log,
		givenName: name,
		fields:    fields,
	}

	ctx = context.WithValue(ctx, existingSpan, s)
	return ctx, s
}

func (s *span) AddPair(ctx context.Context, pair kv.Pair) {
	s.fields = append(s.fields, pair)
}

func (s *span) Complete() {
	elems := []interface{}{fmt.Sprintf("[ %s ]", s.givenName)}
	for _, field := range s.fields {
		elems = append(elems, field)
	}

	s.log(elems...)
}

func (s *span) name() string {
	return s.givenName
}
