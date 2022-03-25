// Package oops offers helper methods to create errors.
// These errors can be of a specific type and will contain contextual
// details regarding an error.
package oops

import (
	"errors"
	"fmt"

	"github.com/thisiserico/golib/kv"
)

// With creates a new error, mergin previously key-value pairs with the
// new ones given.
func With(err error, pairs ...kv.Pair) error {
	var details []kv.Pair

	if structured, ok := err.(*structuredError); ok {
		details = structured.details
	}
	details = append(details, pairs...)

	return &structuredError{
		origin:  err,
		details: details,
	}
}

// Details extracts all the key-value pairs from the given error.
func Details(err error) []kv.Pair {
	if err == nil {
		return nil
	}

	var structured *structuredError
	if !errors.As(err, &structured) {
		return Details(errors.Unwrap(err))
	}

	return append(structured.details, Details(errors.Unwrap(structured))...)
}

// Detail will find the key-value pair from the given error if exists, or an
// empty pair otherwise. An indicator for the pair existence is returned as well.
func Detail(err error, key string) (kv.Pair, bool) {
	for _, pair := range Details(err) {
		if pair.Name() == key {
			return pair, true
		}
	}

	return kv.Pair{}, false
}

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
