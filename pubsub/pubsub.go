// Package pubsub is an abstraction of any publisher-subscriber mechanism.
package pubsub

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thisiserico/golib/v2/kv"
)

// ID indicates the event identifier.
type ID string

// Name indicates the event type.
type Name string

// Meta encapsulates execution information.
type Meta struct {
	// CreatedAtUTC indicates the UTC time when the event was created.
	CreatedAtUTC time.Time `json:"created_at_utc"`

	// CorrelationID holds the request correlation ID.
	CorrelationID string `json:"correlation_id"`

	// Attempts indicates how many times the event has been handled.
	Attempts int `json:"attempts"`

	// IsDryRun indicates whether the execution is a dry run.
	IsDryRun bool `json:"is_dry_run"`
}

// Event defines the event envelope.
type Event struct {
	// ID holds the event unique ID, to be used for idempotency purposes.
	ID ID `json:"id"`

	// Name indicates the event name or type.
	Name Name `json:"name"`

	// Meta encapsulates contextual information.
	Meta Meta `json:"meta"`

	// Payload holds the actual event message.
	Payload []byte `json:"payload"`
}

// NewEvent creates an event of the specified name that uses contextual
// information and the given message.
func NewEvent(ctx context.Context, name Name, msg []byte) Event {
	return Event{
		ID:   ID(uuid.New().String()),
		Name: name,
		Meta: Meta{
			CreatedAtUTC:  time.Now().UTC(),
			CorrelationID: kv.CorrelationID(ctx).String(),
			IsDryRun:      kv.IsDryRun(ctx).Bool(),
		},
		Payload: msg,
	}
}

// Publisher defines the capabilities of any publisher.
type Publisher interface {
	// Emit publishes the given events to the stream.
	Emit(context.Context, ...Event) error

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
