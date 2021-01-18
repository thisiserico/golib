// Package redis provides a way to interact with redis streams.
// segmentio/redis-go is used underneath.
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

// Stream let's both, producers and consumers, know what redis streams to
// interact with.
type Stream struct {
	// name specifies the stream name.
	name string

	// events specifies the event names that will use the specific stream in
	// order to publish events.
	events []string
}

// StreamForPublisher needs to be used to indicate what streams a concrete event
// type needs to be sent to.
func StreamForPublisher(name string, events ...string) Stream {
	return Stream{
		name:   name,
		events: events,
	}
}

// StreamForSubscriber needs to be used to indicate from what streams a
// subscriber needs to get events from.
func StreamForSubscriber(name string) Stream {
	return Stream{
		name: name,
	}
}

type publisher struct {
	// client holds an instance of the redis client, which has an internal
	// state that will be reused.
	client *redis.Client

	// streamForEvent indicates the redis stream to use for a specific
	// event name.
	streamForEvent map[string]string
}

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
	ctx, span := o11y.StartSpan(ctx, "emitter")
	defer span.Complete()

	for i, event := range events {
		stream, exists := p.streamForEvent[string(event.Name)]
		if !exists {
			err := errors.New(
				ctx,
				"unknown stream for event",
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
			err = errors.New(
				ctx,
				"redis publisher",
				err,
				errors.Permanent,
			)
			span.AddPair(ctx, kv.New("error", err))

			return err
		}
	}

	return nil
}

func (p *publisher) Close() error {
	return nil
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

// WithReadingSize indicates how many events can be taken out of the stream
// at once. Defaults to 10.
func WithReadingSize(readSize int) SubscriberOption {
	return func(sub *subscriber) {
		sub.readSize = readSize
	}
}

type subscriber struct {
	client     *redis.Client
	groupID    string
	consumerID string
	// streams []string
	streams     []interface{}
	maxAttempts int
	readSize    int
}

func Subscriber(groupID, address string, streams []Stream, opts ...SubscriberOption) pubsub.Subscriber {
	if len(streams) < 1 {
		panic("at least one stream to subscribe to is required")
	}

	var strs []interface{}
	// var strs []string
	for _, stream := range streams {
		strs = append(strs, stream.name)
	}

	sub := &subscriber{
		client: &redis.Client{
			Addr: address,
		},
		groupID:     groupID,
		consumerID:  uuid.New().String(),
		streams:     strs,
		maxAttempts: 1,
		readSize:    10,
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

	for _, stream := range s.streams {
		resp := s.client.Query(ctx, "xautoclaim", stream, s.groupID, s.consumerID, 0, "0-0", "count", s.readSize)
		var nextStartID interface{}
		_ = resp.Next(&nextStartID)

		var entries []interface{}
		_ = resp.Next(&entries)

		for _, redisEntry := range entries {
			s.consumeSingleEntry(
				ctx,
				stream.(string),
				redisEntry,
				handler,
				errHandler,
			)
		}
		if err := resp.Close(); err != nil {
			errHandler(ctx, errors.New(ctx, "redis xautoclaim error", err, errors.Permanent), nil)
		}
	}

	args := []interface{}{"group", s.groupID, s.consumerID, "count", s.readSize, "block", 0, "streams"}
	args = append(args, s.streams...)
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

			for _, redisEntry := range entries {
				s.consumeSingleEntry(
					ctx,
					string(redisStreams[0].([]byte)),
					redisEntry,
					handler,
					errHandler,
				)
			}
		}
		if err := resp.Close(); err != nil {
			errHandler(ctx, errors.New(ctx, "redis xreadgroup error", err, errors.Permanent), nil)
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
	entry := redisEntry.([]interface{})
	entryID := string(entry[0].([]byte))
	fields := entry[1].([]interface{})

	var event pubsub.Event
	_ = json.Unmarshal(fields[1].([]byte), &event)

	for event.Meta.Attempts < s.maxAttempts {
		event.Meta.Attempts++

		if err := handler(ctx, event); err != nil {
			eventForErrorHandler := &event
			if event.Meta.Attempts != s.maxAttempts {
				eventForErrorHandler = nil
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

// createConsumerGroupForEachStream ensures that the consumer group exists for
// all the streams. This allows to consume from all the streams at once using
// a single XREADGROUP command.
func (s *subscriber) createConsumerGroupForEachStream(ctx context.Context) {
	for _, stream := range s.streams {
		_ = s.client.Query(ctx, "xgroup", "create", stream, s.groupID, "$", "mkstream")
	}
}

func (s *subscriber) Close() error {
	return nil
}
