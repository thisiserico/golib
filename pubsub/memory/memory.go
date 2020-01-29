package memory

import (
	"context"

	"github.com/thisiserico/golib/v2/errors"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/pubsub"
)

var subscribers []*subscriber

type publisher struct{}

// PublisherOption allows to tweak publisher behavior while hidding the
// library internals.
type PublisherOption func(*publisher)

// NewPublisher creates a new in memory publisher implementation.
func NewPublisher(_ ...PublisherOption) pubsub.Publisher {
	return &publisher{}
}

// Emit will publish the provided events to all the existing subscribers.
// This is a blocking operation, no errors will be produced.
func (p *publisher) Emit(ctx context.Context, events ...pubsub.Event) error {
	for _, sub := range subscribers {
		sub.emitEvents(events...)
	}

	return nil
}

func (p *publisher) Close() error { return nil }

type subscriber struct {
	maxAttempts int
	events      chan pubsub.Event
}

// SubscriberOption allows to tweak subscriber behavior while hidding the
// library internals.
type SubscriberOption func(*subscriber)

// WithMaxAttempts indicates how many times an event will be processed if the
// handler erroes. Defaults to 1.
func WithMaxAttempts(retries int) SubscriberOption {
	return func(sub *subscriber) {
		sub.maxAttempts = retries
	}
}

// WithQueueSize indicates how many events can be in flight at any given time.
// Defaults to 10.
func WithQueueSize(queueSize int) SubscriberOption {
	return func(sub *subscriber) {
		sub.events = make(chan pubsub.Event, queueSize)
	}
}

// NewSubscriber creates a new in memory subscriber implementation.
func NewSubscriber(opts ...SubscriberOption) pubsub.Subscriber {
	sub := &subscriber{
		maxAttempts: 1,
		events:      make(chan pubsub.Event, 10),
	}

	for _, opt := range opts {
		opt(sub)
	}

	subscribers = append(subscribers, sub)
	return sub
}

func (s *subscriber) emitEvents(events ...pubsub.Event) {
	for _, event := range events {
		s.events <- event
	}
}

// Consume will handle the given event. The error handler will be used if the
// event handler erroes. Retries will take place as indicated, passing along an
// error only when there're still retries left, an error and the actual event
// otherwise. The error will always contain the handling attempt as a tag.
func (s *subscriber) Consume(ctx context.Context, handler pubsub.Handler, errorHandler pubsub.ErrorHandler) {
	select {
	case <-ctx.Done():
		errorHandler(ctx, errors.New(ctx), nil)

	case event := <-s.events:
		for event.Meta.Attempts < s.maxAttempts {
			event.Meta.Attempts += 1

			if err := handler(ctx, event); err != nil {
				eventForErrorHandler := &event
				if event.Meta.Attempts == s.maxAttempts {
					eventForErrorHandler = nil
				}

				errorHandler(
					ctx,
					errors.New(err, kv.New("attempt", event.Meta.Attempts)),
					eventForErrorHandler,
				)
			}
		}
	}
}

func (s *subscriber) Close() error {
	close(s.events)
	return nil
}
