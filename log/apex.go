package log

import (
	"context"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/thisiserico/golib/constant"
	"github.com/thisiserico/golib/contxt"
)

var _ Logger = &logger{}

type logger struct {
	apex *log.Logger
}

// NewLogger obtains a new logger using the specified configuration.
func NewLogger(of outputFormat) Logger {
	var handler log.Handler
	if of == PlainFormat {
		handler = cli.New(os.Stdout)
	} else {
		handler = json.New(os.Stdout)
	}

	return &logger{
		apex: &log.Logger{
			Level:   log.InfoLevel,
			Handler: handler,
		},
	}
}

func (l *logger) Info(ctx context.Context, msg string, tags Tags) {
	l.prepareApexEntry(ctx, tags).Info(msg)
}

func (l *logger) Error(ctx context.Context, err error, tags Tags) {
	l.prepareApexEntry(ctx, tags).Error(err.Error())
}

func (l *logger) Fatal(ctx context.Context, err error, tags Tags) {
	l.prepareApexEntry(ctx, tags).Fatal(err.Error())
}

func (l *logger) prepareApexEntry(ctx context.Context, tags Tags) *log.Entry {
	entry := log.NewEntry(l.apex)
	entry = entry.WithField(string(constant.BuildID), contxt.BuildID(ctx))
	entry = entry.WithField(string(constant.CorrelationID), contxt.CorrelationID(ctx))
	entry = entry.WithField(string(constant.InBehalfOf), contxt.InBehalfOfService(ctx))
	entry = entry.WithField(string(constant.IsDryRun), contxt.IsDryRunExecution(ctx))
	entry = entry.WithField(string(constant.ServiceHost), contxt.RunningInHost(ctx))
	entry = entry.WithField(string(constant.ServiceName), contxt.RunningService(ctx))
	entry = entry.WithField(string(constant.WhosRequesting), contxt.RequestedByService(ctx))

	for key, value := range tags {
		entry = entry.WithField(string(key), value)
	}

	return entry
}
