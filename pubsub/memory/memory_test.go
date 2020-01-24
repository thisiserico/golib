package memory

import (
	"context"
	"testing"

	"github.com/thisiserico/golib/errors"
	"github.com/thisiserico/golib/pubsub"
)

func TestTheConsumerWithACanceledContext(t *testing.T) {
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
	sub.Consume(ctx, handler, errHandler)

	if messageWasHandled {
		t.Fatal("no messages should have been handled")
	}

	if !errors.Is(obtainedError, errors.ContextError) {
		t.Fatalf("a context error was expected, got %v", obtainedError)
	}
}

func TestWithSingleSubscription(t *testing.T) {
	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		return nil
	}

	errHandler := func(_ context.Context, _ error, _ *pubsub.Event) {}

	pub := NewPublisher()
	sub := NewSubscriber()

	event := pubsub.NewEvent(context.Background(), pubsub.Name("type"), nil)
	_ = pub.Emit(context.Background(), event)

	sub.Consume(context.Background(), handler, errHandler)

	if !messageWasHandled {
		t.Fatal("the messages should have been handled")
	}
}

func TestWithMultipleSubscriptions(t *testing.T) {
	var messageWasHandledOnSubscriber1 bool
	handlerForSubscriber1 := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandledOnSubscriber1 = true
		return nil
	}

	var messageWasHandledOnSubscriber2 bool
	handlerForSubscriber2 := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandledOnSubscriber2 = true
		return nil
	}

	errHandler := func(_ context.Context, _ error, _ *pubsub.Event) {}

	pub := NewPublisher()
	sub1 := NewSubscriber()
	sub2 := NewSubscriber()

	event := pubsub.NewEvent(context.Background(), pubsub.Name("type"), nil)
	_ = pub.Emit(context.Background(), event)

	sub1.Consume(context.Background(), handlerForSubscriber1, errHandler)
	sub2.Consume(context.Background(), handlerForSubscriber2, errHandler)

	if !messageWasHandledOnSubscriber1 {
		t.Fatal("the messages should have been handled")
	}
	if !messageWasHandledOnSubscriber2 {
		t.Fatal("the messages should have been handled")
	}
}

func TestWhenTheHandlerFails(t *testing.T) {
	var messageWasHandled bool
	handler := func(_ context.Context, _ pubsub.Event) error {
		messageWasHandled = true
		return errors.New("handler error")
	}

	var errorWasHandled bool
	errHandler := func(_ context.Context, _ error, _ *pubsub.Event) {
		errorWasHandled = true
	}

	pub := NewPublisher()
	sub := NewSubscriber()

	event := pubsub.NewEvent(context.Background(), pubsub.Name("type"), nil)
	_ = pub.Emit(context.Background(), event)

	sub.Consume(context.Background(), handler, errHandler)

	if !messageWasHandled {
		t.Fatal("the messages should have been handled")
	}
	if !errorWasHandled {
		t.Fatal("the error should have been handled")
	}
}
