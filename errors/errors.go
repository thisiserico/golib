// Package errors provides a way to generate contextual errors.
package errors

import (
	"context"
	"strings"

	"github.com/thisiserico/golib/constant"
	"github.com/thisiserico/golib/contxt"
)

const (
	// Contextual indicates that the error contains contextual tags.
	Contextual = Type(iota)

	// ContextError indicates a context.Context error.
	ContextError

	// Permanent indicates that the error is permanent – use along Transient.
	Permanent

	// PlainError indicates that the error is not contextual.
	PlainError

	// Transient indicates that the error is transient – use along Permanent.
	Transient
)

// Type indicates the error type on its inception.
type Type int

// TagPair encapsulates a key value pair.
type TagPair struct {
	key   constant.Key
	value constant.Value
}

// Tag pairs a key with a value.
func Tag(k constant.Key, v constant.Value) TagPair {
	return TagPair{
		key:   k,
		value: v,
	}
}

var _ error = contextualError{}

type contextualError struct {
	types []Type
	msgs  []string
	tags  map[constant.Key]constant.Value
}

// Error satisfies the error contract.
func (err contextualError) Error() string {
	return strings.Join(err.msgs, ": ")
}

// New facilitates the contextual error creation by accepting different argument
// types: context, error, message, type and tags.
//
// - nil
//   Getting a nil `nil` value explicetely means that there was no error.
// - context.Context
//   If the context has erroed, a ContextError type is added to the list of
//   types. The tags will be populated with the known contextual values and
//   Contextual type.
// - string
//   The given message is added to the error stack.
// - contextualError
//   To keep the consistency between errors.
// - error
//   The error message is added to the error stack.
// - Type
//   The given type is stacked. It can be later be accessed with `Is`.
// - TagPair
//   To provide more contextual data points.
func New(args ...interface{}) error {
	types := make([]Type, 0)
	msgs := make([]string, 0)
	tags := make(map[constant.Key]constant.Value)

	for _, arg := range args {
		switch t := arg.(type) {
		case nil:
			return nil

		case context.Context:
			select {
			case <-t.Done():
				types = append(types, ContextError)

			default:
			}

			types = append(types, Contextual)
			tags[constant.BuildID] = contxt.BuildID(t)
			tags[constant.ServiceHost] = contxt.RunningInHost(t)
			tags[constant.ServiceName] = contxt.RunningService(t)
			tags[constant.CorrelationID] = contxt.CorrelationID(t)
			tags[constant.InBehalfOf] = contxt.InBehalfOfService(t)
			tags[constant.WhosRequesting] = contxt.RequestedByService(t)
			tags[constant.IsDryRun] = contxt.IsDryRunExecution(t)

		case string:
			msgs = append(msgs, t)

		case contextualError:
			types = append(types, t.types...)
			msgs = append(msgs, t.msgs...)
			tags = t.tags

		case error:
			msgs = append(msgs, t.Error())

		case Type:
			types = append(types, t)

		case TagPair:
			if _, exists := tags[t.key]; !exists {
				tags[t.key] = t.value
			}
		}
	}

	return contextualError{
		types: types,
		msgs:  msgs,
		tags:  tags,
	}
}

// Is returns true when the given error stack contains the requested type.
func Is(anyError error, requested Type) bool {
	err, isContextual := anyError.(contextualError)
	if !isContextual {
		return false
	}

	for _, t := range err.types {
		if t == requested {
			return true
		}
	}

	return requested == PlainError
}

// Tags returns a key-value dictionary.
func Tags(anyError error) map[constant.Key]constant.Value {
	err, isContextual := anyError.(contextualError)
	if !isContextual {
		return nil
	}

	return err.tags
}
