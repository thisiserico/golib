// Package memory provides an agent implementation to use when testing
// observability elements.
//
// The agent registers itself against the o11y package so that it becomes
// transparent for the test what's being used underneath.
package memory

import (
	"context"

	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/o11y"
)

const (
	defaultSpanName = "new"
	existingSpan    = contextKey("existing_span")
)

var ag *agent

func init() {
	if ag != nil {
		return
	}

	ag = &agent{}
	o11y.Register(ag)
}

type contextKey string

type agent struct{}

// StartSpan acts a bit different than expected, as it will only create a new
// span when one doesn't yet exist in the context. This allows for the test to
// create the span itself and let production code use such span directly.
func (a *agent) StartSpan(ctx context.Context, name string) (context.Context, o11y.Span) {
	if val := ctx.Value(existingSpan); val != nil {
		return ctx, val.(*span)
	}

	return newSpan(ctx, defaultSpanName)
}

// GetSpan, on the other hand, works as you'd expect.
func (a *agent) GetSpan(ctx context.Context) o11y.Span {
	_, s := a.StartSpan(ctx, defaultSpanName)
	return s
}

func (a *agent) Flush() {}

type span struct {
	name      string
	fields    []kv.Pair
	err       error
	completed bool
}

func newSpan(ctx context.Context, name string) (context.Context, o11y.Span) {
	s := &span{
		name:   name,
		fields: make([]kv.Pair, 0),
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
	s.completed = true
}

// IsCompleted informs whether a provided span has been completed yet.
func IsCompleted(s o11y.Span) bool {
	return s.(*span).completed
}

// HasErroed indicates whether an error has been reported into the span.
func HasErroed(s o11y.Span) bool {
	return s.(*span).err != nil
}

// PairMatches indicates whether a pair with the provided key and value has
// been reported.
func PairMatches(s o11y.Span, key string, val interface{}) bool {
	var matches bool
	for _, pair := range s.(*span).fields {
		if pair.Name() != key {
			continue
		}

		// Do not return to match against the last occurrence of the pair,
		// not the first one found.
		matches = pair.Value() == val
	}

	return matches
}
