package memory

import (
	"context"
	"testing"
	"time"

	"github.com/thisiserico/golib/v2/errors"
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

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sub := NewSubscriber()
	defer sub.Close()
	sub.Consume(ctx, handler, errHandler)

	if messageWasHandled {
		t.Fatal("no messages should have been handled")
	}
	if !errors.Is(obtainedError, errors.Context) {
		t.Fatalf("a context error was expected, got %v", obtainedError)
	}
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
	defer sub.Close()

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	sub.Consume(context.Background(), handler, errHandler)

	if obtainedError == nil {
		t.Fatal("an error had to be handled")
	}

	pair, exists := errors.Tag("attempt", obtainedError)
	if !exists {
		t.Fatal("the handling attempt had to be present in the error")
	}
	if got := pair.Int(); got != 1 {
		t.Fatalf("invalid handling attempt, want 1, got %d", got)
	}
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
	defer sub.Close()

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	sub.Consume(context.Background(), handler, errHandler)

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

	pair, exists := errors.Tag("attempt", obtainedErrors[1])
	if !exists {
		t.Fatal("the handling attempt had to be present in the error")
	}
	if got := pair.Int(); got != maxAttempts {
		t.Fatalf("invalid handling attempt, want %d, got %d", maxAttempts, got)
	}
}

func TestASubscriberWithAFilledUpQueue(t *testing.T) {
	pub := NewPublisher()
	sub := NewSubscriber(WithQueueSize(1))
	defer sub.Close()

	event := pubsub.NewEvent(context.Background(), knownEventName, nil)
	_ = pub.Emit(context.Background(), event)

	var lastEventEmitted bool
	go func() {
		_ = pub.Emit(context.Background(), event)
		lastEventEmitted = true
	}()

	<-time.After(100 * time.Millisecond)
	if lastEventEmitted {
		t.Fatal("the second event shouldn't be emitted")
	}
}

func TestUsingMultipleSubscribers(t *testing.T) {
	var handledEvents []pubsub.Event
	handler := func(_ context.Context, e pubsub.Event) error {
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

	sub1.Consume(context.Background(), handler, errHandler)
	sub2.Consume(context.Background(), handler, errHandler)

	if got := len(handledEvents); got != 2 {
		t.Fatalf("the same event had to be handled twice, it's been handled %d", got)
	}

	ev1 := handledEvents[0].ID
	ev2 := handledEvents[1].ID
	if ev1 != ev2 {
		t.Fatal("the handled events don't match")
	}
}
