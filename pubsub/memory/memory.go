// Package memory provides an in-memory pubsub mechanism.
// Its usage is recommended when operating a monolith, an actual pubsub engine
// has not yet been chosen or when working on POCs and the like.
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/o11y"
	"github.com/thisiserico/golib/oops"
	"github.com/thisiserico/golib/pubsub"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	lock        sync.RWMutex
	subscribers map[string]*subscriber
)

func init() {
	lock.Lock()
	defer lock.Unlock()

	subscribers = make(map[string]*subscriber)
}

type publisher struct {
	id     string
	tracer trace.Tracer
}

// PublisherOption allows to tweak publisher behavior while hidding the
// library internals.
type PublisherOption func(*publisher)

// WithPublisherID indicates the publisher identifier for traceability purposes.
func WithPublisherID(id string) PublisherOption {
	return func(pub *publisher) {
		pub.id = id
	}
}

// NewPublisher creates a new in memory publisher implementation.
func NewPublisher(opts ...PublisherOption) pubsub.Publisher {
	pub := &publisher{
		id:     uuid.New().String(),
		tracer: otel.Tracer("pubsub/memory.publisher"),
	}

	for _, opt := range opts {
		opt(pub)
	}

	return pub
}

// Emit will publish the provided events to all the existing subscribers.
// This is a blocking operation, no errors will be produced.
func (p *publisher) Emit(ctx context.Context, events ...pubsub.Event) error {
	ctx, span := p.tracer.Start(
		ctx,
		"emit",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attribute.String("pubsub.publisher_id", p.id)),
		trace.WithAttributes(o11y.Attributes(ctx)...),
	)
	defer span.End()

	for _, ev := range events {
		span.AddEvent(string(ev.Name))
	}

	lock.RLock()
	defer lock.RUnlock()
	for _, sub := range subscribers {
		sub.emitEvents(events...)
	}

	return nil
}

func (p *publisher) Close() error { return nil }

type subscriber struct {
	id          string
	maxAttempts int
	events      chan pubsub.Event
	tracer      trace.Tracer
}

// SubscriberOption allows to tweak subscriber behavior while hidding the
// library internals.
type SubscriberOption func(*subscriber)

// WithSubscriberID indicates the consumer identifier for traceability purposes.
func WithSubscriberID(id string) SubscriberOption {
	return func(sub *subscriber) {
		sub.id = id
	}
}

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
		close(sub.events)
		sub.events = make(chan pubsub.Event, queueSize)
	}
}

// NewSubscriber creates a new in memory subscriber implementation.
func NewSubscriber(opts ...SubscriberOption) pubsub.Subscriber {
	sub := &subscriber{
		id:          uuid.New().String(),
		maxAttempts: 1,
		events:      make(chan pubsub.Event, 10),
		tracer:      otel.Tracer("pubsub/memory.subscriber"),
	}

	for _, opt := range opts {
		opt(sub)
	}

	lock.Lock()
	defer lock.Unlock()
	subscribers[sub.id] = sub

	return sub
}

func (s *subscriber) emitEvents(events ...pubsub.Event) {
	defer func() {
		// After closing a subscriber events channel, pending events can cause
		// a panic.
		recover()
	}()

	for _, event := range events {
		s.events <- event
	}
}

// Consume will consume events as they are available. The error handler will be
// used if the event handler erroes. Retries will take place as indicated,
// passing along an error only when there're still retries left, an error and
// the actual event otherwise. The error will always contain the handling
// attempt as a tag.
func (s *subscriber) Consume(ctx context.Context, handler pubsub.Handler, errorHandler pubsub.ErrorHandler) {
	for {
		if err := ctx.Err(); err != nil {
			break
		}

		s.consumeEvent(ctx, handler, errorHandler)
	}
}

func (s *subscriber) consumeEvent(ctx context.Context, handler pubsub.Handler, errorHandler pubsub.ErrorHandler) {
	select {
	case <-ctx.Done():
		return

	case event := <-s.events:
		ctx := kv.SetDynamicAttributes(ctx, event.Meta.CorrelationID, event.Meta.IsDryRun)
		ctx, span := s.tracer.Start(
			ctx,
			"consume",
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String("pubsub.subscriber_id", s.id),
				attribute.String("pubsub.event_name", string(event.Name)),
			),
			trace.WithAttributes(o11y.Attributes(ctx)...),
		)
		defer span.End()

		for event.Meta.Attempts < s.maxAttempts {
			span.AddEvent(fmt.Sprintf("attempt %d", event.Meta.Attempts))
			event.Meta.Attempts++

			err := handler(ctx, event)
			if err == nil {
				break
			}
			span.RecordError(err)

			eventForErrorHandler := &event
			if event.Meta.Attempts != s.maxAttempts {
				eventForErrorHandler = nil
			} else {
				span.SetStatus(codes.Error, "max attempts reached")
			}

			errorHandler(
				ctx,
				oops.With(
					err,
					kv.New("pubsub.attempt", event.Meta.Attempts),
					kv.New("pubsub.is_last_attempt", event.Meta.Attempts == s.maxAttempts),
				),
				eventForErrorHandler,
			)
		}
	}
}

func (s *subscriber) Close() error {
	lock.Lock()
	defer lock.Unlock()

	close(s.events)
	delete(subscribers, s.id)

	return nil
}
