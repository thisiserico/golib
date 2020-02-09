// Package halt exposes a convenience method to deal with graceful shutdowns.
package halt

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger"
)

// Halter will be used to wait for shutdown requests.
type Halter interface {
	// Wait should block until a shutdown is requested.
	Wait()
}

var _ Halter = &halter{}

type halter struct {
	ctx context.Context
	log logger.Log
}

// New configures and returns the context to use when shutting down.
func New(ctx context.Context, log logger.Log) (context.Context, Halter) {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		log(ctx, "stopper signal captured", kv.New("signal", <-stop))
		cancel()
	}()

	return ctx, &halter{
		ctx: ctx,
		log: log,
	}
}

func (s *halter) Wait() {
	<-s.ctx.Done()
	s.log(s.ctx, "stopper gracefully shuting down")
}
