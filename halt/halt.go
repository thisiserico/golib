// Package halt exposes a convenience method to deal with grafecul shutdowns.
package halt

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thisiserico/golib/constant"
	"github.com/thisiserico/golib/log"
)

// Halter will be used to wait for shutdown requests.
type Halter interface {
	// Wait should block until a shutdown is requested.
	Wait()
}

var _ Halter = &halter{}

type halter struct {
	ctx    context.Context
	logger log.Logger
}

// New configures and returns the context.Context to use when shutting down.
func New(ctx context.Context, logger log.Logger) (context.Context, Halter) {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		logger.Info(ctx, "stopper signal captured", log.Tags{constant.Key("signal"): constant.AnyValue(<-stop)})
		cancel()
	}()

	return ctx, &halter{
		ctx:    ctx,
		logger: logger,
	}
}

func (s *halter) Wait() {
	<-s.ctx.Done()
	s.logger.Info(s.ctx, "stopper gracefully shuting down", nil)
}
