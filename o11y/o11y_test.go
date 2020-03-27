package o11y

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/thisiserico/golib/v2/kv"
)

func TestThatEverythingIsDelegatedToTheAgent(t *testing.T) {
	t.Run("when starting a span", func(t *testing.T) {
		ag := newTestingAgent()
		Register(ag)

		name := uuid.New().String()
		_, s := StartSpan(context.Background(), name)

		if got := len(ag.startedSpans); got != 1 {
			t.Fatalf("one span had to be created, created %d instead", got)
		}
		if span := ag.startedSpans[0]; span != s {
			t.Fatalf("the span was unexpected, want %#v, got %#v", s, span)
		}
	})

	t.Run("when obtaining a span from the context", func(t *testing.T) {
		ag := newTestingAgent()
		Register(ag)

		name := uuid.New().String()
		ctx, originalSpan := StartSpan(context.Background(), name)
		newSpan := GetSpan(ctx)

		if got := len(ag.startedSpans); got != 1 {
			t.Fatalf("one span had to be created, created %d instead", got)
		}
		if originalSpan != newSpan {
			t.Fatalf("the span was unexpected, want %#v, got %#v", originalSpan, newSpan)
		}
	})

	t.Run("when flusing the agent", func(t *testing.T) {
		ag := newTestingAgent()
		Register(ag)

		ag.Flush()
		if !ag.flushed {
			t.Fatal("the agent had to be flushed")
		}
	})
}

type testingAgent struct {
	flushed      bool
	startedSpans []*testingSpan
}

type testingSpan struct {
	name string
}

type ctxKey string

func newTestingAgent() *testingAgent {
	return &testingAgent{
		startedSpans: make([]*testingSpan, 0),
	}
}

func newTestingSpan(name string) *testingSpan {
	return &testingSpan{
		name: name,
	}
}

func (t *testingAgent) StartSpan(ctx context.Context, name string) (context.Context, Span) {
	span := newTestingSpan(name)
	t.startedSpans = append(t.startedSpans, span)

	ctx = context.WithValue(ctx, ctxKey("span"), span)

	return ctx, span
}

func (t *testingAgent) GetSpan(ctx context.Context) Span {
	return ctx.Value(ctxKey("span")).(*testingSpan)
}

func (t *testingAgent) Flush() {
	t.flushed = true
}

func (t *testingSpan) AddPair(_ context.Context, _ kv.Pair) {}

func (t *testingSpan) Complete() {}
