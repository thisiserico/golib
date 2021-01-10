package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/thisiserico/golib/v2/errors"
	"github.com/thisiserico/golib/v2/o11y"
	"github.com/thisiserico/golib/v2/o11y/memory"
	"github.com/thisiserico/golib/v2/pubsub"
)

var (
	knownEventName = pubsub.Name("known")
	errHandler     = errors.New("handler error")
)

func TestConsumingWithACanceledContext(t *testing.T) {
	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		return nil
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		obtainedError = err
	}

	ctx, _ := o11y.StartSpan(context.Background(), "")
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	sub := NewSubscriber()
	sub.Consume(ctx, handler, errHandler)

	if messageWasHandled {
		t.Fatal("no messages should have been handled")
	}
	if obtainedError != nil {
		t.Fatalf("no errors were expected, got %v", obtainedError)
	}
	if memory.IsCompleted(o11y.GetSpan(ctx)) {
		t.Fatal("no span should have been created, therefore completed")
	}

	sub.Close()
}

func TestAHandlerThatFails(t *testing.T) {
	handler := func(_ context.Context, _ pubsub.Event) error {
		return errHandler
	}

	var obtainedError error
	errHandler := func(_ context.Context, err error, _ *pubsub.Event) {
		obtainedError = err
	}

	pub := NewPublisher()
	sub := NewSubscriber()

	pubCtx, _ := o11y.StartSpan(context.Background(), "")
	event1 := pubsub.NewEvent(pubCtx, knownEventName, nil)
	event2 := pubsub.NewEvent(pubCtx, knownEventName, nil)
	_ = pub.Emit(pubCtx, event1, event2)

	if !memory.IsCompleted(o11y.GetSpan(pubCtx)) {
		t.Fatal("the publisher span should have been completed")
	}
	if memory.HasErroed(o11y.GetSpan(pubCtx)) {
		t.Fatal("no errors were expected after publishing the event")
	}
	if !memory.PairMatches(o11y.GetSpan(pubCtx), "event_0", knownEventName) {
		t.Fatal("the event_0 attribute should be reported")
	}
	if !memory.PairMatches(o11y.GetSpan(pubCtx), "event_1", knownEventName) {
		t.Fatal("the event_1 attribute should be reported")
	}

	subCtx, _ := o11y.StartSpan(context.Background(), "")
	subCtx, _ = context.WithTimeout(subCtx, 50*time.Millisecond)
	sub.Consume(subCtx, handler, errHandler)

	if obtainedError == nil {
		t.Fatal("an error had to be handled")
	}
	if !memory.IsCompleted(o11y.GetSpan(subCtx)) {
		t.Fatal("the publisher span should have been completed")
	}
	if !memory.HasErroed(o11y.GetSpan(subCtx)) {
		t.Fatal("an error was expected after consuming an event")
	}
	if !memory.PairMatches(o11y.GetSpan(subCtx), "attempt", 0) {
		t.Fatal("the number of attemps were not reported correctly")
	}

	pair, exists := errors.Tag("attempt", obtainedError)
	if !exists {
		t.Fatal("the handling attempt had to be present in the error")
	}
	if got := pair.Int(); got != 1 {
		t.Fatalf("invalid handling attempt, want 1, got %d", got)
	}

	sub.Close()
}

func TestAHandlerThatFailsMultipleTimes(t *testing.T) {
	const maxAttempts = 2

	handler := func(_ context.Context, _ pubsub.Event) error {
		return errHandler
	}

	var (
		obtainedErrors []error
		obtainedEvents []*pubsub.Event
	)
	errHandler := func(_ context.Context, err error, event *pubsub.Event) {
		obtainedErrors = append(obtainedErrors, err)
		obtainedEvents = append(obtainedEvents, event)
	}

	pub := NewPublisher()
	sub := NewSubscriber(WithMaxAttempts(maxAttempts))

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	ctx, _ := o11y.StartSpan(context.Background(), "")
	ctx, _ = context.WithTimeout(ctx, 50*time.Millisecond)
	sub.Consume(ctx, handler, errHandler)

	if got := len(obtainedErrors); got != maxAttempts {
		t.Fatalf("as many errors as handling attempts are expected, want %d, got %d", maxAttempts, got)
	}

	if !memory.PairMatches(o11y.GetSpan(ctx), "attempt", 1) {
		t.Fatal("the number of attemps were not reported correctly")
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

	pair, exists := errors.Tag("is_last_attempt", obtainedErrors[0])
	if !exists {
		t.Fatal("the is_last_attempt indicator had to be present in the error")
	}
	if got := pair.Bool(); got != false {
		t.Fatalf("invalid is_last_attempt, want %t, got %T", false, got)
	}

	pair, exists = errors.Tag("attempt", obtainedErrors[1])
	if !exists {
		t.Fatal("the handling attempt had to be present in the error")
	}
	if got := pair.Int(); got != maxAttempts {
		t.Fatalf("invalid handling attempt, want %d, got %d", maxAttempts, got)
	}

	pair, exists = errors.Tag("is_last_attempt", obtainedErrors[1])
	if !exists {
		t.Fatal("the is_last_attempt indicator had to be present in the error")
	}
	if got := pair.Bool(); got != true {
		t.Fatalf("invalid is_last_attempt, want %t, got %T", true, got)
	}

	sub.Close()
}

func TestASubscriberWithAFilledUpQueue(t *testing.T) {
	pub := NewPublisher()
	sub := NewSubscriber(WithQueueSize(1))

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	var lastEventEmitted bool
	go func() {
		_ = pub.Emit(context.Background(), event)
		lastEventEmitted = true
	}()

	<-time.After(50 * time.Millisecond)
	if lastEventEmitted {
		t.Fatal("the second event shouldn't be emitted")
	}

	handler := func(_ context.Context, _ pubsub.Event) error { return nil }
	errHandler := func(_ context.Context, _ error, _ *pubsub.Event) {}

	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	sub.Consume(ctx, handler, errHandler)
	sub.Close()
}

func TestUsingMultipleSubscribers(t *testing.T) {
	var lock sync.Mutex
	var handledEvents []pubsub.Event
	handler := func(_ context.Context, e pubsub.Event) error {
		lock.Lock()
		defer lock.Unlock()

		handledEvents = append(handledEvents, e)
		return nil
	}

	errHandler := func(_ context.Context, _ error, _ *pubsub.Event) {}

	pub := NewPublisher()
	sub1 := NewSubscriber()
	sub2 := NewSubscriber()
	defer sub1.Close()
	defer sub2.Close()

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	go sub1.Consume(context.Background(), handler, errHandler)
	go sub2.Consume(context.Background(), handler, errHandler)
	<-time.After(50 * time.Millisecond)

	lock.Lock()
	defer lock.Unlock()
	if got := len(handledEvents); got != 2 {
		t.Fatalf("the same event had to be handled twice, it's been handled %d times", got)
	}

	ev1 := handledEvents[0].ID
	ev2 := handledEvents[1].ID
	if ev1 != ev2 {
		t.Fatal("the handled events don't match")
	}
}
