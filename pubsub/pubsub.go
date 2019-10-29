// Package pubsub is an abstraction of any publisher-subscriber mechanism.
package pubsub

import (
	"context"
	"time"

	"github.com/thisiserico/golib/constant"
)

// Name indicates the event type.
type Name string

// Meta encapsulates execution information.
type Meta struct {
	createdAt     time.Time
	correlationID constant.OneCorrelationID  `json:"correlation_id"`
	isDryRun      constant.IsDryRunExecution `json:"is_dry_run"`
}

// Event defines the event envelope.
type Event struct {
	id   constant.OneCorrelationID `json:"id"`
	name Name                      `json:"name"`
	msg  interface{}               `json:"msg"`
	meta Meta                      `json:"meta"`
}

// Publisher defines the capabilities of any publisher.
type Publisher interface {
	// Emit publishes the given event to the stream.
	Emit(context.Context, Event) error

	// Close should close any underlying connection.
	Close() error
}

// Subscriber defines the capabilities of any subscriber.
type Subscriber interface {
	// Consume takes a single event from the stream and processes it using the
	// given handler. The event will be sent to the error handler in case of
	// failure.
	Consume(context.Context, Handler, ErrorHandler)

	// Close should close any underlying connection.
	Close() error
}

// Handler handles the given event.
type Handler func(context.Context, Event) error

// ErrorHandler should handle the given error and deal with the optional
// event if necessary.
type ErrorHandler func(context.Context, error, *Event)
