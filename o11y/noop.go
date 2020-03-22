package o11y

import (
	"context"

	"github.com/thisiserico/golib/v2/kv"
)

func init() {
	if agent != nil {
		return
	}

	agent = noop()
}

type noopAgent struct{}

func noop() Agent {
	return &noopAgent{}
}

func (a *noopAgent) StartSpan(ctx context.Context, _ string) (context.Context, Span) {
	return newSpan(ctx)
}

func (a *noopAgent) GetSpan(ctx context.Context) Span {
	_, s := newSpan(ctx)
	return s
}

func (a *noopAgent) Flush() {}

type span struct{}

func newSpan(ctx context.Context) (context.Context, Span) {
	return ctx, &span{}
}

func (s *span) AddPair(ctx context.Context, pair kv.Pair) {}

func (s *span) Complete() {}
