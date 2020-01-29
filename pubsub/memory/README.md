# memory
--
    import "github.com/thisiserico/golib/pubsub/memory"


## Usage

#### func  NewPublisher

```go
func NewPublisher(_ ...PublisherOption) pubsub.Publisher
```
NewPublisher creates a new in memory publisher implementation.

#### func  NewSubscriber

```go
func NewSubscriber(opts ...SubscriberOption) pubsub.Subscriber
```
NewSubscriber creates a new in memory subscriber implementation.

#### type PublisherOption

```go
type PublisherOption func(*publisher)
```

PublisherOption allows to tweak publisher behavior while hidding the library
internals.

#### type SubscriberOption

```go
type SubscriberOption func(*subscriber)
```

SubscriberOption allows to tweak subscriber behavior while hidding the library
internals.

#### func  WithMaxAttempts

```go
func WithMaxAttempts(retries int) SubscriberOption
```
WithMaxAttempts indicates how many times an event will be processed if the
handler erroes. Defaults to 1.

#### func  WithQueueSize

```go
func WithQueueSize(queueSize int) SubscriberOption
```
WithQueueSize indicates how many events can be in flight at any given time.
Defaults to 10.
