# pubsub
--
    import "github.com/thisiserico/golib/pubsub"

Package pubsub is an abstraction of any publisher-subscriber mechanism.

## Usage

#### type ErrorHandler

```go
type ErrorHandler func(context.Context, error, *Event)
```

ErrorHandler should handle the given error and deal with the optional event if
necessary.

#### type Event

```go
type Event struct {
	// ID holds the event unique ID, to be used for dempotency purposes.
	ID ID `json:"id"`

	// Name indicates the event name or type.
	Name Name `json:"name"`

	// Meta encapsulates contextual information.
	Meta Meta `json:"meta"`

	// Payload holds the actual event message.
	Payload interface{} `json:"payload"`
}
```

Event defines the event envelope.

#### func  NewEvent

```go
func NewEvent(ctx context.Context, name Name, msg interface{}) Event
```
NewEvent creates an event of the specified name that uses contextual information
and the given message.

#### type Handler

```go
type Handler func(context.Context, Event) error
```

Handler handles the given event.

#### type ID

```go
type ID string
```

ID indicates the event identifier.

#### type Meta

```go
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
```

Meta encapsulates execution information.

#### type Name

```go
type Name string
```

Name indicates the event type.

#### type Publisher

```go
type Publisher interface {
	// Emit publishes the given events to the stream.
	Emit(context.Context, ...Event) error

	// Close should close any underlying connection.
	Close() error
}
```

Publisher defines the capabilities of any publisher.

#### type Subscriber

```go
type Subscriber interface {
	// Consume takes a single event from the stream and processes it using the
	// given handler. The event will be sent to the error handler in case of
	// failure.
	Consume(context.Context, Handler, ErrorHandler)

	// Close should close any underlying connection.
	Close() error
}
```

Subscriber defines the capabilities of any subscriber.
