// Package o11y provides an abstraction to make system observable.
// Although the unit of work is the span, other nuances are being dealt with
// underneath, like logging when necessary.
package o11y

import (
	"context"

	"github.com/thisiserico/golib/v2/kv"
)

var agent Agent

// Register allows clients to specify the agent they want to use under the
// hood. All operations will be delegated to a specific agent. Calling this
// method more than once will overwrite the previously specified strategy.
// By default, a no-op agent is used.
func Register(a Agent) {
	agent = a
}

// Agent defines a way to either materialize an existing span in the context
// or to create a new one from scratch.
type Agent interface {
	// StartSpan needs to be used when a new portion of work just started.
	StartSpan(context.Context, string) (context.Context, Span)

	// GetSpan on the other hand, is encouraged to be used when the program
	// is still dealing with the same unit of work.
	GetSpan(context.Context) Span

	// Flush completes any pending span in the program.
	Flush()
}

// Span defines a way to inform the system about the context of a certain
// execution.
type Span interface {
	// AddPair allows to provide such context.
	AddPair(context.Context, kv.Pair)

	// Complete is used when the unit of work has completed.
	Complete()
}

// StartSpan generates a new span using the specified agent. Each agent will
// have to provide the span in the context.
func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return agent.StartSpan(ctx, name)
}

// GetSpan on the other hand, tries to extract an existing span from the given
// context. Each agent will deal with the fact that a previous span might not
// exist yet.
func GetSpan(ctx context.Context) Span {
	return agent.GetSpan(ctx)
}
