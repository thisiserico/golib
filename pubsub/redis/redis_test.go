// +build redis

package redis

import (
	"context"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/redis-go"
	"github.com/thisiserico/golib/v2/errors"
	"github.com/thisiserico/golib/v2/pubsub"
)

var redisAddress = flag.String("address", "127.0.0.1:6379", "redis string and port (defaults to 127.0.0.1:6379)")

func TestThatUnknownEventsCannotBeEmitted(t *testing.T) {
	pub := Publisher(*redisAddress, nil)
	event := pubsub.NewEvent(context.Background(), "unknown_event_name", nil)

	err := pub.Emit(context.Background(), event)
	if err == nil {
		t.Fatal("an error was expected, got none")
	}
	if !errors.Is(err, errors.NonExistent) {
		t.Fatalf("a non existent error was expected, got %#v", err)
	}
	if !errors.Is(err, errors.Permanent) {
		t.Fatalf("a permanent error was expected, got %#v", err)
	}
}

func TestThatNothingGetsEmittedWhenTheContextIsCancelled(t *testing.T) {
	stream := uuid.New().String()
	eventName := uuid.New().String()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pub.Emit(ctx, event)
	if err == nil {
		t.Fatal("an error was expected, got none")
	}
	if !errors.Is(err, errors.Context) {
		t.Fatalf("a non existent error was expected, got %#v", err)
	}
	if !errors.Is(err, errors.Permanent) {
		t.Fatalf("a permanent error was expected, got %#v", err)
	}
}

func TestThatASubscriberRequiresAtLeastOneStream(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("a panic was expected, got none")
		}
	}()

	_ = Subscriber(uuid.New().String(), *redisAddress, nil)
}

func TestThatNothingGetsConsumedWhenTheContextIsCancelled(t *testing.T) {
	groupID := uuid.New().String()
	stream := uuid.New().String()

	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		if errors.Is(err, errors.Context) {
			return
		}
		obtainedError = err
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sub := Subscriber(groupID, *redisAddress, StreamsForSubscriber(stream))
	sub.Consume(ctx, handler, errHandler)

	if messageWasHandled {
		t.Fatal("no message should have been handled")
	}
	if obtainedError != nil {
		t.Fatalf("no error was expected, got %#v", obtainedError)
	}
}

func TestThatNothingGetsConsumedWhenNoPreviousEventsHaveBeenEmitted(t *testing.T) {
	groupID := uuid.New().String()
	stream := uuid.New().String()

	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		if errors.Is(err, errors.Context) {
			return
		}
		obtainedError = err
	}

	sub := Subscriber(groupID, *redisAddress, StreamsForSubscriber(stream))
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	go sub.Consume(ctx, handler, errHandler)

	<-ctx.Done()
	if messageWasHandled {
		t.Fatal("no message should have been handled")
	}
	if obtainedError != nil {
		t.Fatalf("no error was expected, got %#v", obtainedError)
	}
}

func TestThatAFailedHandlingReportsAnError(t *testing.T) {
	groupID := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	handler := func(_ context.Context, _ pubsub.Event) error {
		return errors.New("handler error")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		obtainedError = err
		cancel()
	}

	sub := Subscriber(groupID, *redisAddress, StreamsForSubscriber(stream))
	go sub.Consume(ctx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctx.Done()
	if obtainedError == nil {
		t.Fatal("an error had to be handled")
	}
}

func TestThatHandlingAttempsCanBeRequestedOnFailedHandlings(t *testing.T) {
	const maxAttempts = 2

	groupID := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	handler := func(_ context.Context, event pubsub.Event) error {
		return errors.New("handler error")
	}

	ctx, cancel := context.WithCancel(context.Background())

	var obtainedErrors []error
	var obtainedEvents []*pubsub.Event
	errHandler := func(_ context.Context, err error, event *pubsub.Event) {
		obtainedErrors = append(obtainedErrors, err)
		obtainedEvents = append(obtainedEvents, event)

		if len(obtainedErrors) == maxAttempts {
			cancel()
		}
	}

	sub := Subscriber(
		groupID,
		*redisAddress,
		StreamsForSubscriber(stream),
		HandlingNumberOfAttempts(maxAttempts),
	)
	go sub.Consume(ctx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctx.Done()
	if got := len(obtainedErrors); got != maxAttempts {
		t.Fatalf("as many errors as handling attempts are expected, want %d, got %d", maxAttempts, got)
	}

	if got := obtainedEvents[0]; got != nil {
		t.Fatalf("non-final handling attempts should not report the event, got %#v", got)
	}
	if obtainedEvents[1] == nil {
		t.Fatal("the last handling attempt should report the event")
	}
	if obtainedEvents[1].ID != event.ID {
		t.Fatal("the reported event doesn't match the expected one")
	}
}

func TestConsumingOneEventFromOneStream(t *testing.T) {
	groupID := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())

	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		cancel()

		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		obtainedError = err
	}

	sub := Subscriber(groupID, *redisAddress, StreamsForSubscriber(stream))
	go sub.Consume(ctx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctx.Done()
	if !messageWasHandled {
		t.Fatal("a message had to be handled")
	}
	if obtainedError != nil {
		t.Fatalf("no handling errors were expected, got %#v", obtainedError)
	}
}

func TestConsumingMultipleEventsFromOneStream(t *testing.T) {
	const readSize = 2

	groupID := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())

	var obtainedEvents []pubsub.Event
	handler := func(_ context.Context, event pubsub.Event) error {
		obtainedEvents = append(obtainedEvents, event)
		if len(obtainedEvents) == readSize {
			cancel()
		}

		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		if errors.Is(err, errors.Context) {
			return
		}
		obtainedError = err
	}

	sub := Subscriber(
		groupID,
		*redisAddress,
		StreamsForSubscriber(stream),
		ReadingBatchCapacity(readSize),
	)

	// Trigger initial consumption so that the redis stream and consumer group
	// exist, positioning the cursor in 0-0 for the group in the stream.
	// This allows to later on read from ">" successfully.
	setupCtx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	sub.Consume(setupCtx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event, event)

	go sub.Consume(ctx, handler, errHandler)

	<-ctx.Done()
	if got := len(obtainedEvents); got != readSize {
		t.Fatalf("as many successful handlings as events were expected, want %d, got %d", readSize, got)
	}
	if obtainedError != nil {
		t.Fatalf("no handling errors were expected, got %#v", obtainedError)
	}
}

func TestConsumingOneEventFromMultipleStreams(t *testing.T) {
	groupID := uuid.New().String()
	streamOne := uuid.New().String()
	streamTwo := uuid.New().String()
	eventName := uuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())

	var messageWasHandled bool
	handler := func(_ context.Context, event pubsub.Event) error {
		messageWasHandled = true
		cancel()

		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		obtainedError = err
	}

	sub := Subscriber(
		groupID,
		*redisAddress,
		StreamsForSubscriber(streamOne, streamTwo),
	)
	go sub.Consume(ctx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(streamOne, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctx.Done()
	if !messageWasHandled {
		t.Fatal("a message had to be handled")
	}
	if obtainedError != nil {
		t.Fatalf("no handling errors were expected, got %#v", obtainedError)
	}
}

func TestConsumingMultipleEventsFromMultipleStreams(t *testing.T) {
	const readSize = 2

	groupID := uuid.New().String()
	streamOne := uuid.New().String()
	streamTwo := uuid.New().String()
	eventNameOne := uuid.New().String()
	eventNameTwo := uuid.New().String()

	ctx, cancel := context.WithCancel(context.Background())

	var obtainedEvents []pubsub.Event
	handler := func(_ context.Context, event pubsub.Event) error {
		obtainedEvents = append(obtainedEvents, event)
		if len(obtainedEvents) == readSize {
			cancel()
		}

		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		if errors.Is(err, errors.Context) {
			return
		}

		obtainedError = err
	}

	sub := Subscriber(
		groupID,
		*redisAddress,
		StreamsForSubscriber(streamOne, streamTwo),
		ReadingBatchCapacity(readSize),
	)

	// Trigger initial consumption so that the redis stream and consumer group
	// exist, positioning the cursor in 0-0 for the group in the stream.
	// This allows to later on read from ">" successfully.
	setupCtx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	sub.Consume(setupCtx, handler, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{
		StreamForPublisher(streamOne, eventNameOne),
		StreamForPublisher(streamTwo, eventNameTwo),
	})
	eventOne := pubsub.NewEvent(context.Background(), pubsub.Name(eventNameOne), nil)
	eventTwo := pubsub.NewEvent(context.Background(), pubsub.Name(eventNameTwo), nil)
	_ = pub.Emit(context.Background(), eventOne, eventTwo)

	go sub.Consume(ctx, handler, errHandler)

	<-ctx.Done()
	if got := len(obtainedEvents); got != readSize {
		t.Fatalf("as many successful handlings as events were expected, want %d, got %d", readSize, got)
	}
	if obtainedError != nil {
		t.Fatalf("no handling errors were expected, got %#v", obtainedError)
	}
}

func TestConsumingFromMultipleSubscribers(t *testing.T) {
	groupIDOne := uuid.New().String()
	groupIDTwo := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	ctxOne, cancelOne := context.WithCancel(context.Background())
	ctxTwo, cancelTwo := context.WithCancel(context.Background())

	var handlerOneHandledMessage bool
	handlerOne := func(_ context.Context, _ pubsub.Event) error {
		handlerOneHandledMessage = true
		cancelOne()

		return nil
	}
	var handlerTwoHandledMessage bool
	handlerTwo := func(_ context.Context, _ pubsub.Event) error {
		handlerTwoHandledMessage = true
		cancelTwo()

		return nil
	}
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {}

	subOne := Subscriber(groupIDOne, *redisAddress, StreamsForSubscriber(stream))
	subTwo := Subscriber(groupIDTwo, *redisAddress, StreamsForSubscriber(stream))
	go subOne.Consume(ctxOne, handlerOne, errHandler)
	go subTwo.Consume(ctxTwo, handlerTwo, errHandler)
	leaveTimeForTheSubscriberToStartRunning()

	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctxOne.Done()
	<-ctxTwo.Done()
	if !handlerOneHandledMessage {
		t.Fatal("a message had to be handled from the first subscriber")
	}
	if !handlerTwoHandledMessage {
		t.Fatal("a message had to be handled from the second subscriber")
	}
}

func TestThatStreamEntriesAreNeverLost(t *testing.T) {
	const expectedEvents = 3

	groupID := uuid.New().String()
	stream := uuid.New().String()
	eventName := uuid.New().String()

	event := pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	js, _ := json.Marshal(event)

	client := &redis.Client{Addr: *redisAddress}
	// Make sure a consumer group exists for the stream.
	_ = client.Query(context.Background(), "xgroup", "create", stream, groupID, "$", "mkstream")
	// An event is added to and read from the stream, but never acknowledged, making it claimed.
	_ = client.Exec(context.Background(), "xadd", stream, "*", "event", js)
	_ = client.Query(context.Background(), "xreadgroup", "group", groupID, uuid.New().String(), "count", 1, "streams", stream, ">")

	// Another event makes it to the stream before it's being consumed, making it not yet claimed.
	pub := Publisher(*redisAddress, []Stream{StreamForPublisher(stream, eventName)})
	event = pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	ctx, cancel := context.WithCancel(context.Background())

	var obtainedEvents []pubsub.Event
	handler := func(_ context.Context, event pubsub.Event) error {
		obtainedEvents = append(obtainedEvents, event)
		if len(obtainedEvents) == expectedEvents {
			cancel()
		}

		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		if errors.Is(err, errors.Context) {
			return
		}
		obtainedError = err
	}

	sub := Subscriber(groupID, *redisAddress, StreamsForSubscriber(stream))
	go sub.Consume(ctx, handler, errHandler)

	// Yet another event is produced, this one while the subscriber is already consuming.
	event = pubsub.NewEvent(context.Background(), pubsub.Name(eventName), nil)
	_ = pub.Emit(context.Background(), event)

	<-ctx.Done()
	if got := len(obtainedEvents); got != expectedEvents {
		t.Fatalf("as many successful handlings as events were expected, want %d, got %d", expectedEvents, got)
	}
	if obtainedError != nil {
		t.Fatalf("no handling errors were expected, got %#v", obtainedError)
	}
}

func leaveTimeForTheSubscriberToStartRunning() {
	<-time.After(100 * time.Millisecond)
}
