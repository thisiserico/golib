// Package redis provides a way to interact with redis streams.
// segmentio/redis-go is used underneath.
// Redis version >= 6.2 is required.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/redis-go"
	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/o11y"
	"github.com/thisiserico/golib/oops"
	"github.com/thisiserico/golib/pubsub"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const streamCapacity = 1000000

// Stream let's both –producers and consumers– know what redis streams to interact with.
type Stream struct {
	// name holds the stream name.
	name string

	// events holds the event names that will use the this stream in order to publish events.
	events []string

	// capacity indicates how many approximate events will live in a stream at most. This
	// attribute can be read as XTRIM name MAXLEN ~ capacity. Defaults to streamCapacity.
	capacity int
}

// CappedAt let clients indicate the maximum capacity for this specific redis stream. Notice that,
// in the scenario of capping the same stream twice with different values, this client won't do
// any special handling: it will overwrite previous capacity references.
func (s Stream) CappedAt(capacity int) Stream {
	s.capacity = streamCapacity
	return s
}

// StreamForPublisher creates a Stream that will know the event types that will use a redis stream
// with such name in order to publish messages into.
func StreamForPublisher(name string, events ...string) Stream {
	return Stream{
		name:     name,
		events:   events,
		capacity: streamCapacity,
	}
}

// StreamsForSubscriber creates a list of Stream to indicate all the redis streams a subscriber
// has to read events from.
func StreamsForSubscriber(names ...string) []Stream {
	streams := make([]Stream, 0, len(names))
	for _, name := range names {
		streams = append(streams, Stream{name: name})
	}

	return streams
}

type publisher struct {
	// client holds an instance of the redis client, which has an internal
	// state that will be reused.
	client *redis.Client

	// streamForEvent indicates the redis stream to use for a specific
	// event name.
	streamForEvent map[string]string

	// streamCapacities indicates the capacity of each redis stream.
	streamCapacities map[string]int

	tracer trace.Tracer
}

// Publisher creates a publisher that uses redis streams under the hood.
func Publisher(address string, streams []Stream) pubsub.Publisher {
	streamForEvents := make(map[string]string)
	streamCapacities := make(map[string]int)
	for _, stream := range streams {
		for _, event := range stream.events {
			streamForEvents[event] = stream.name
		}

		streamCapacities[stream.name] = stream.capacity
	}

	return &publisher{
		client: &redis.Client{
			Addr: address,
		},
		streamForEvent:   streamForEvents,
		streamCapacities: streamCapacities,
		tracer:           otel.Tracer("pubsub/redis.publisher"),
	}
}

func (p *publisher) Emit(ctx context.Context, events ...pubsub.Event) error {
	ctx, span := p.tracer.Start(
		ctx,
		"emit",
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(o11y.Attributes(ctx)...),
	)
	defer span.End()

	for _, event := range events {
		stream, exists := p.streamForEvent[string(event.Name)]
		if !exists {
			err := oops.With(
				oops.NonExistent("unknown redis stream for event"),
				kv.New("event_name", event.Name),
			)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}
		span.AddEvent(string(event.Name))

		js, _ := json.Marshal(event)
		redisCapacity := p.streamCapacities[stream]
		err := p.client.Exec(ctx, "xadd", stream, "maxlen", redisCapacity, "*", "event", js)
		if err != nil {
			err = oops.Invalid("redis xadd: %w", err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return err
		}
	}

	return nil
}

func (p *publisher) Close() error {
	return nil
}

// SubscriberOption allows to tweak subscriber behavior.
type SubscriberOption func(*subscriber)

// HandlingNumberOfAttempts indicates how many times an event will be processed if the handler
// errors. Defaults to 1, that is, no automatic retries.
func HandlingNumberOfAttempts(attempts int) SubscriberOption {
	return func(sub *subscriber) {
		sub.maxAttempts = attempts
	}
}

// ReadingBatchCapacity indicates how many events can be taken out of the stream at once.
// Defaults to 10.
func ReadingBatchCapacity(capacity int) SubscriberOption {
	return func(sub *subscriber) {
		sub.readCapacity = capacity
	}
}

// ConsumeTimeout indicates the maximum amount of time for an event to be in a handling state.
// Defaults to 1s, which is the minimum value.
func ConsumeTimeout(timeout time.Duration) SubscriberOption {
	return func(sub *subscriber) {
		if timeout.Seconds() < 1 {
			timeout = time.Second
		}

		sub.consumeTimeout = timeout
	}
}

// RunFailureRecovery enables the execution of the redis xautoclaim command, running it on the
// indicated cadence. By default, no recovery is run.
func RunFailureRecovery(enabled bool, cadence time.Duration) SubscriberOption {
	return func(sub *subscriber) {
		sub.failureRecoveryEnabled = enabled
		sub.failureRecoveryCadence = cadence
	}
}

type subscriber struct {
	client                 *redis.Client
	groupID                string
	consumerID             string
	streams                []string
	maxAttempts            int
	readCapacity           int
	consumeTimeout         time.Duration
	failureRecoveryEnabled bool
	failureRecoveryCadence time.Duration
	tracer                 trace.Tracer
}

// Subscriber creates a subscriber that uses redis streams under the hood.
// All the events that are handled (either successfully or by using the error handler), won't be
// consumed again. On the other hand, only events that can't be handled by the client will be
// re-consumed automatically.
// This makes the error handler responsible for dealing with unsuccessful handlings. The use of
// DLQs is encouraged to ensure all events are processed.
func Subscriber(groupID, address string, streams []Stream, opts ...SubscriberOption) pubsub.Subscriber {
	if len(streams) < 1 {
		panic("at least one stream to read from is required")
	}

	var strs []string
	for _, stream := range streams {
		strs = append(strs, stream.name)
	}

	sub := &subscriber{
		client: &redis.Client{
			Addr: address,
		},
		groupID:                groupID,
		consumerID:             uuid.New().String(),
		streams:                strs,
		maxAttempts:            1,
		readCapacity:           10,
		consumeTimeout:         time.Second,
		failureRecoveryEnabled: false,
		failureRecoveryCadence: time.Second,
		tracer:                 otel.Tracer("pubsub/redis.subscriber"),
	}

	for _, opt := range opts {
		opt(sub)
	}

	return sub
}

func (s *subscriber) Consume(
	ctx context.Context,
	handler pubsub.Handler,
	errHandler pubsub.ErrorHandler,
) {
	s.createConsumerGroupForEachStream(ctx)
	s.handleClaimedButNotProcessedEvents(ctx, handler, errHandler)

	// There has to be an easier way to compose the list below...
	args := []interface{}{"group", s.groupID, s.consumerID, "count", s.readCapacity, "block", 0, "streams"}
	for _, stream := range s.streams {
		args = append(args, stream)
	}
	for range s.streams {
		args = append(args, ">")
	}

	for {
		if err := ctx.Err(); err != nil {
			break
		}

		resp := s.client.Query(ctx, "xreadgroup", args...)
		var redisResponse interface{}
		for resp.Next(&redisResponse) {
			redisStreams := redisResponse.([]interface{})
			entries := redisStreams[1].([]interface{})
			streamID := string(redisStreams[0].([]byte))

			for _, redisEntry := range entries {
				s.consumeSingleEntry(ctx, streamID, redisEntry, handler, errHandler)
			}
		}
		if err := resp.Close(); err != nil {
			err = oops.Invalid("redis xreadgroup: %w", err)
			errHandler(ctx, err, nil)
		}
	}
}

// createConsumerGroupForEachStream ensures that the consumer group exists for
// all the streams. This allows to consume from all the streams at once using
// a single XREADGROUP command.
func (s *subscriber) createConsumerGroupForEachStream(ctx context.Context) {
	ctx, span := s.tracer.Start(
		ctx,
		"consumer group set up",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("pubsub.group_id", s.groupID),
			attribute.String("pubsub.consumer_id", s.consumerID),
		),
		trace.WithAttributes(o11y.Attributes(ctx)...),
	)
	defer span.End()

	for _, stream := range s.streams {
		span.AddEvent(stream)
		_ = s.client.Query(ctx, "xgroup", "create", stream, s.groupID, "$", "mkstream")
	}
}

func (s *subscriber) handleClaimedButNotProcessedEvents(
	ctx context.Context,
	handler pubsub.Handler,
	errHandler pubsub.ErrorHandler,
) {
	idleTimeout := time.Duration(s.maxAttempts) * s.consumeTimeout

	go func() {
		for ; s.failureRecoveryEnabled; <-time.After(s.failureRecoveryCadence) {
			ctx, span := s.tracer.Start(
				ctx,
				"potential failure recovery",
				trace.WithSpanKind(trace.SpanKindConsumer),
				trace.WithAttributes(attribute.String("pubsub.group_id", s.groupID)),
				trace.WithAttributes(o11y.Attributes(ctx)...),
			)

			for _, stream := range s.streams {
				span.AddEvent(stream)

				resp := s.client.Query(ctx, "xautoclaim", stream, s.groupID, s.consumerID, idleTimeout, "0-0", "count", s.readCapacity)
				var nextStartID interface{}
				_ = resp.Next(&nextStartID)

				var entries []interface{}
				_ = resp.Next(&entries)

				for _, redisEntry := range entries {
					s.consumeSingleEntry(ctx, stream, redisEntry, handler, errHandler)
				}
				if err := resp.Close(); err != nil {
					err = oops.Invalid("redis xautoclaim: %w", err)
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())

					errHandler(ctx, err, nil)
				}
			}

			span.End()
		}
	}()
}

// consumeSingleEntry handles a single redis entry, acknowledging the entry at
// the end, no matter whether it was successfully handled or not. This makes
// the error handler responsible to handle errors in any way fits.
func (s *subscriber) consumeSingleEntry(
	ctx context.Context,
	streamID string,
	redisEntry interface{},
	handler pubsub.Handler,
	errHandler pubsub.ErrorHandler,
) {
	ctx, span := s.tracer.Start(
		ctx,
		"consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("pubsub.group_id", s.groupID),
			attribute.String("pubsub.consumer_id", s.consumerID),
		),
		trace.WithAttributes(o11y.Attributes(ctx)...),
	)
	defer span.End()

	ctx, _ = context.WithTimeout(ctx, s.consumeTimeout)

	entry := redisEntry.([]interface{})
	entryID := string(entry[0].([]byte))
	fields := entry[1].([]interface{})

	var event pubsub.Event
	_ = json.Unmarshal(fields[1].([]byte), &event)

	ctx = kv.SetDynamicAttributes(ctx, event.Meta.CorrelationID, event.Meta.IsDryRun)
	span.SetAttributes(
		attribute.String("pubsub.correlation_id", event.Meta.CorrelationID),
		attribute.Bool("pubsub.is_dry_run", event.Meta.IsDryRun),
	)
	span.AddEvent(string(event.Name))

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
			span.SetStatus(codes.Error, err.Error())
		}

		errHandler(
			ctx,
			oops.With(
				err,
				kv.New("pubsub.attempt", event.Meta.Attempts),
				kv.New("pubsub.is_last_attempt", event.Meta.Attempts == s.maxAttempts),
			),
			eventForErrorHandler,
		)
	}

	_ = s.client.Exec(context.Background(), "xack", streamID, s.groupID, entryID)
}

func (s *subscriber) Close() error {
	return nil
}
