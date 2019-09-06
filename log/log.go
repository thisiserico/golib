package log

import (
	"context"

	"github.com/thisiserico/golib/constant"
)

const (
	// PlainFormat uses plain text as output.
	PlainFormat outputFormat = iota

	// JSONFormat uses JSON as output (used by default).
	JSONFormat
)

type outputFormat int

// Tags abstracts a key -> value dictionary.
type Tags map[constant.Key]constant.Value

//go:generate mockgen -package=log -destination=./mock.go github.com/thisiserico/golib/log Logger

// Logger defines the used log capabilities.
type Logger interface {
	// Info specifies an informative log entry.
	Info(context.Context, string, Tags)

	// Error specifies an error log entry.
	Error(context.Context, error, Tags)

	// Fatal specifies a fatal log entry.
	Fatal(context.Context, error, Tags)
}
