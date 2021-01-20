// Package redis provides a way to interact with redis streams.
// segmentio/redis-go is used underneath.
// Redis version >= 6.2 is required.
package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/redis-go"
	"github.com/thisiserico/golib/v2/errors"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/o11y"
	"github.com/thisiserico/golib/v2/pubsub"
)

// Stream let's both –producers and consumers– know what redis streams to interact with.
type Stream struct {
	// name holds the stream name.
	name string

	// events holds the event names that will use the this stream in order to publish events.
	events []string
}

// StreamForPublisher creates a Stream that will know the event types that will use a redis stream
// with such name in order to publish messages into.
func StreamForPublisher(name string, events ...string) Stream {
	return Stream{
		name:   name,
		events: events,
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
}

// Publisher creates a publisher that uses redis streams under the hood.
func Publisher(address string, streams []Stream) pubsub.Publisher {
	streamForEvents := make(map[string]string)
	for _, stream := range streams {
		for _, event := range stream.events {
			streamForEvents[event] = stream.name
		}
	}

	return &publisher{
		client: &redis.Client{
			Addr: address,
		},
		streamForEvent: streamForEvents,
	}
}

func (p *publisher) Emit(ctx context.Context, events ...pubsub.Event) error {
	ctx, span := o11y.StartSpan(ctx, "redis emitter")
	defer span.Complete()

	for i, event := range events {
		stream, exists := p.streamForEvent[string(event.Name)]
		if !exists {
			err := errors.New(
				ctx,
				"unknown redis stream for event",
				kv.New("event_name", event.Name),
				errors.NonExistent,
				errors.Permanent,
			)
			span.AddPair(ctx, kv.New("error", err))

			return err
		}
		span.AddPair(ctx, kv.New(fmt.Sprintf("event_%d", i), event.Name))

		js, _ := json.Marshal(event)
		err := p.client.Exec(ctx, "xadd", stream, "*", "event", js)
		if err != nil {
			err = errors.New(ctx, "redis xadd", err, errors.Permanent)
			span.AddPair(ctx, kv.New("error", err))

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

type subscriber struct {
	client       *redis.Client
	groupID      string
	consumerID   string
	streams      []string
	maxAttempts  int
	readCapacity int
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
		groupID:      groupID,
		consumerID:   uuid.New().String(),
		streams:      strs,
		maxAttempts:  1,
		readCapacity: 10,
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
			err = errors.New(ctx, "redis xreadgroup", err, errors.Permanent)
			errHandler(ctx, err, nil)
		}
	}
}

// createConsumerGroupForEachStream ensures that the consumer group exists for
// all the streams. This allows to consume from all the streams at once using
// a single XREADGROUP command.
func (s *subscriber) createConsumerGroupForEachStream(ctx context.Context) {
	ctx, span := o11y.StartSpan(ctx, "redis consumer group set up")
	defer span.Complete()

	span.AddPair(ctx, kv.New("group_id", s.groupID))
	span.AddPair(ctx, kv.New("consumer_id", s.consumerID))

	for _, stream := range s.streams {
		_ = s.client.Query(ctx, "xgroup", "create", stream, s.groupID, "$", "mkstream")
	}
}

func (s *subscriber) handleClaimedButNotProcessedEvents(
	ctx context.Context,
	handler pubsub.Handler,
	errHandler pubsub.ErrorHandler,
) {
	// TODO run this asynchronously, using a s.consumeTimeout property as "idle-time".

	ctx, span := o11y.StartSpan(ctx, "redis potential failure recovery")
	defer span.Complete()

	span.AddPair(ctx, kv.New("group_id", s.groupID))
	span.AddPair(ctx, kv.New("consumer_id", s.consumerID))

	for _, stream := range s.streams {
		resp := s.client.Query(ctx, "xautoclaim", stream, s.groupID, s.consumerID, 0, "0-0", "count", s.readCapacity)
		var nextStartID interface{}
		_ = resp.Next(&nextStartID)

		var entries []interface{}
		_ = resp.Next(&entries)

		for _, redisEntry := range entries {
			s.consumeSingleEntry(ctx, stream, redisEntry, handler, errHandler)
		}
		if err := resp.Close(); err != nil {
			err = errors.New(ctx, "redis xautoclaim", err, errors.Permanent)
			span.AddPair(ctx, kv.New("error", err))
			errHandler(ctx, err, nil)
		}
	}
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
	ctx, span := o11y.StartSpan(ctx, "redis consumer")
	defer span.Complete()

	span.AddPair(ctx, kv.New("group_id", s.groupID))
	span.AddPair(ctx, kv.New("consumer_id", s.consumerID))

	entry := redisEntry.([]interface{})
	entryID := string(entry[0].([]byte))
	fields := entry[1].([]interface{})

	var event pubsub.Event
	_ = json.Unmarshal(fields[1].([]byte), &event)

	span.AddPair(ctx, kv.New("event_name", event.Name))

	for event.Meta.Attempts < s.maxAttempts {
		span.AddPair(ctx, kv.New("attempt", event.Meta.Attempts))
		event.Meta.Attempts++

		if err := handler(ctx, event); err != nil {
			eventForErrorHandler := &event
			if event.Meta.Attempts != s.maxAttempts {
				eventForErrorHandler = nil
			} else {
				span.AddPair(ctx, kv.New("error", err))
			}

			errHandler(
				ctx,
				errors.New(
					err,
					kv.New("attempt", event.Meta.Attempts),
					kv.New("is_last_attempt", event.Meta.Attempts == s.maxAttempts),
				),
				eventForErrorHandler,
			)
		}
	}

	_ = s.client.Exec(context.Background(), "xack", streamID, s.groupID, entryID)
}

func (s *subscriber) Close() error {
	return nil
}
