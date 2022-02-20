// Package logger provides a simplified logger to write log lines into a writer.
package logger

import (
	"context"
	"io"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/oops"
)

const (
	// PlainOutput writes plain text into the log writer.
	PlainOutput Output = iota

	// JSONOutput writes json encoded text into the log writer.
	JSONOutput
)

// Output defines a way to format log lines before sending them to the writer.
type Output int

// Log lets clients output log lines into the previously specified writer.
// Different argument types –listed below– will be used to define the log line
// composition. The order is important, as arguments can override. By default,
// info log lines are provided.
//
// - `context.Context`
//   Known execution indicators are extracted from the context and provided in
//   the log line as tags.
// - `string`
//   The argument will be used as the log message.
// - `error`
//   The error message will be used as the log message. An error log line will
//   be provided. Error details will be extracted and used as tags.
// - `kv.Pair`
//   Each pair will be used as a log line tag.
//
// Other types will be ignored.
type Log func(...interface{})

// New provides a new logging method. When used, the output will be sent to the
// indicated writer, previously formatting the log line using the specified
// output method.
func New(w io.Writer, o Output) Log {
	var handler log.Handler
	if o == PlainOutput {
		handler = cli.New(w)
	} else {
		handler = json.New(w)
	}

	l := &logger{
		apex: &log.Logger{
			Level:   log.InfoLevel,
			Handler: handler,
		},
	}

	return l.log
}

type logger struct {
	apex *log.Logger
}

func (l *logger) log(args ...interface{}) {
	var isError bool
	entry := log.NewEntry(l.apex)

	var msg string
	for _, arg := range args {
		switch t := arg.(type) {
		case context.Context:
			for _, attr := range kv.AllAttributes(t) {
				val := attr.Value()
				if val == nil {
					continue
				}

				entry = entry.WithField(attr.Name(), val)
			}

		case string:
			msg = t

		case error:
			isError = true
			msg = t.Error()

			for _, pair := range oops.Details(t) {
				entry = entry.WithField(pair.Name(), pair.Value())
			}

		case kv.Pair:
			entry = entry.WithField(t.Name(), t.Value())
		}
	}

	log := entry.Info
	if isError {
		log = entry.Error
	}

	log(msg)
}
