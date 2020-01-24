package memory

import (
	"context"

	"github.com/thisiserico/golib/errors"
	"github.com/thisiserico/golib/pubsub"
)

var subscribers []*subscriber

type publisher struct{}

// NewPublisher creates a new in memory publisher implementation.
func NewPublisher() pubsub.Publisher {
	return &publisher{}
}

func (p *publisher) Emit(ctx context.Context, e pubsub.Event) error {
	for _, sub := range subscribers {
		sub.emitEvent(e)
	}

	return nil
}

type subscriber struct {
	events chan pubsub.Event
}

// NewSubscriber creates a new in memory subscriber implementation.
func NewSubscriber() pubsub.Subscriber {
	sub := &subscriber{
		events: make(chan pubsub.Event, 10),
	}

	subscribers = append(subscribers, sub)
	return sub
}

func (s *subscriber) emitEvent(ev pubsub.Event) {
	s.events <- ev
}

func (s *subscriber) Consume(ctx context.Context, handler pubsub.Handler, errorHandler pubsub.ErrorHandler) {
	select {
	case <-ctx.Done():
		errorHandler(ctx, errors.New(ctx.Err(), errors.ContextError), nil)

	case ev := <-s.events:
		if err := handler(ctx, ev); err != nil {
			errorHandler(ctx, err, &ev)
		}
	}
}

func (p *publisher) Close() error  { return nil }
func (s *subscriber) Close() error { return nil }
