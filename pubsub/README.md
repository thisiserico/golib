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
}
```

Event defines the event envelope.

#### type Handler

```go
type Handler func(context.Context, Event) error
```

Handler handles the given event.

#### type Meta

```go
type Meta struct {
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
