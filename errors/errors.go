// Package errors provides a way to generate contextual errors.
package errors

import (
	"context"
	"strings"

	"github.com/thisiserico/golib/v2/kv"
)

const (
	// Context indicates a context.Context error.
	Context = Category("context-error")

	// Decode indicates that decoding failed.
	Decode = Category("decode")

	// Encode indicates that encoding failed.
	Encode = Category("encode")

	// Existent indicates that an element already exists.
	Existent = Category("existent")

	// Invalid indicates a validation constraint.
	Invalid = Category("invalid")

	// Permanent indicates that the error is permanent – use along Transient.
	Permanent = Category("permanent")

	// NonExistent indicates that an element doesn't exist.
	NonExistent = Category("non-existent")

	// Transient indicates that the error is transient – use along Permanent.
	Transient = Category("transient")
)

// Category indicates the error type on its inception.
type Category string

var _ error = contextualError{}

type contextualError struct {
	msgs       []string
	categories []Category
	tags       []kv.Pair
}

// New creates a contextual error by accepting different arguments listed
// below. The arguments need to be passed in that order to end up with a
// consistent error.
//
// - `context.Context`
//   If the context has erroed, a Context category is added to the list of
//   categories.
// - `string`
//   The given message is added to the error stack.
// - `contextualError`
//   The existing error stack, categories and tags are preserved and/or extended.
// - `error`
//   The error message is added to the error stack.
// - `Category`
//   The given type is stacked. It can be later be accessed with `Is`.
// - `kv.Pair`
//   To provide more contextual data points.
func New(args ...interface{}) error {
	msgs := make([]string, 0)
	categories := make([]Category, 0)
	tags := make([]kv.Pair, 0)

	for _, arg := range args {
		switch t := arg.(type) {
		case context.Context:
			if err := t.Err(); err != nil {
				categories = append(categories, Context)
			}
			for _, attr := range kv.AllAttributes(t) {
				tags = append(tags, kv.New(attr.Name(), attr.Value))
			}

		case string:
			msgs = append(msgs, t)

		case contextualError:
			msgs = append(msgs, t.msgs...)
			categories = append(categories, t.categories...)
			tags = append(tags, t.tags...)

		case error:
			msgs = append(msgs, t.Error())

		case Category:
			categories = append(categories, t)

		case kv.Pair:
			tags = append(tags, t)
		}
	}

	return contextualError{
		msgs:       msgs,
		categories: categories,
		tags:       tags,
	}
}

// Error satisfies the error contract.
func (err contextualError) Error() string {
	return strings.Join(err.msgs, ": ")
}

// Is evaluates whether the given error matches a particular category.
func Is(err error, cat Category) bool {
	contextualErr, isContextual := err.(contextualError)
	if !isContextual {
		return false
	}

	for _, category := range contextualErr.categories {
		if category == cat {
			return true
		}
	}

	return false
}

// Tags returns the error contextual data points.
func Tags(err error) []kv.Pair {
	contextualErr, isContextual := err.(contextualError)
	if !isContextual {
		return nil
	}

	return contextualErr.tags
}

// Tag returns the requested tag if exists, a nil one otherwise. A boolean
// will indicate whether the tag exists.
func Tag(key string, err error) (kv.Pair, bool) {
	contextualErr, isContextual := err.(contextualError)
	if !isContextual {
		return kv.New(key, nil), false
	}

	for _, tag := range contextualErr.tags {
		if tag.Name() == key {
			return tag, true
		}
	}

	return kv.New(key, nil), false
}
