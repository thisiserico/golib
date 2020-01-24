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
	ID constant.OneCorrelationID `json:"id"`

	// Name indicates the event name or type.
	Name Name `json:"name"`

	// Msg holds the actual event payload.
	Msg interface{} `json:"msg"`

	// Meta encapsulates contextual information.
	Meta Meta `json:"meta"`
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

#### type Meta

```go
type Meta struct {
	// CreatedAtUTC indicates the UTC time when the event was created.
	CreatedAtUTC time.Time `json:"created_at_utc"`

	// CorrelationID holds the request correlation ID.
	CorrelationID constant.OneCorrelationID `json:"correlation_id"`

	// IsDryRun indicates whether the execution is a dry run.
	IsDryRun constant.IsDryRunExecution `json:"is_dry_run"`
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
	// Emit publishes the given event to the stream.
	Emit(context.Context, Event) error

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
