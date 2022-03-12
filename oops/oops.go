// Package oops offers helper methods to create errors.
// These errors can be of a specific type and will contain contextual
// details regarding an error.
package oops

import (
	"errors"
	"fmt"

	"github.com/thisiserico/golib/kv"
)

type structuredError struct {
	typology error
	origin   error
	details  []kv.Pair
}

func newError(err error, msg string, args ...interface{}) error {
	return &structuredError{
		typology: err,
		origin:   fmt.Errorf(msg, args...),
	}
}

func (se structuredError) Error() string {
	return se.origin.Error()
}

func (se *structuredError) Is(target error) bool {
	return errors.Is(se.typology, target) || errors.Is(se.origin, target)
}

func (se structuredError) Unwrap() error {
	return errors.Unwrap(se.origin)
}
