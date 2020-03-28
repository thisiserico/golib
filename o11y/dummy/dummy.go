package dummy

import (
	"context"
	"fmt"

	"github.com/thisiserico/golib/v2/cntxt"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger"
	"github.com/thisiserico/golib/v2/o11y"
)

const (
	defaultSpanName = "new"
	existingSpan    = contextKey("existing_span")
)

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
	log          logger.Log
	hasParent    bool
	child        *span
	wasCompleted bool

	name   string
	fields []kv.Pair
	err    error
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
		log:    log,
		name:   name,
		fields: fields,
	}

	if val := ctx.Value(existingSpan); val != nil {
		s.hasParent = true
		val.(*span).child = s
	}

	ctx = context.WithValue(ctx, existingSpan, s)
	return ctx, s
}

func (s *span) AddPair(ctx context.Context, pair kv.Pair) {
	if err, isError := pair.Value().(error); isError {
		s.err = err
		return
	}

	s.fields = append(s.fields, pair)
}

func (s *span) Complete() {
	if s.wasCompleted {
		return
	}
	if s.hasParent {
		return
	}

	s.wasCompleted = true
	s.dump()
}

func (s *span) dump() {
	elems := []interface{}{fmt.Sprintf("%s", s.name)}

	if s.err != nil {
		elems = append(elems, s.err)
	}

	for _, field := range s.fields {
		elems = append(elems, field)
	}

	s.log(elems...)

	if s.child == nil {
		return
	}
	s.child.dump()
}
