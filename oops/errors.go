package oops

import (
	"errors"
)

var (
	// ErrCancelled indicates that the operation was cancelled.
	ErrCancelled = errors.New("cancelled")

	// ErrDecode indicates that decoding failed.
	ErrDecode = errors.New("decode")

	// ErrEncode indicates that encoding failed.
	ErrEncode = errors.New("encode")

	// ErrExistent indicates that a resource already exists.
	ErrExistent = errors.New("existent")

	// ErrInvalid indicates an unmet validation constraint.
	ErrInvalid = errors.New("invalid")

	// ErrNonExistent indicates that a resource doesn't exist.
	ErrNonExistent = errors.New("non-existent")

	// ErrTimeout indicates that the operation was not fulfiled in time.
	ErrTimeout = errors.New("timeout")

	// ErrTransient suggests that the operation might work if retried.
	ErrTransient = errors.New("transient")
)

// Cancelled creates a new error that satisfies errors.Is(err, ErrCancelled).
func Cancelled(msg string, args ...interface{}) error {
	return newError(ErrCancelled, msg, args...)
}

// Decode creates a new error that satisfies errors.Is(err, ErrDecode).
func Decode(msg string, args ...interface{}) error {
	return newError(ErrDecode, msg, args...)
}

// Encode creates a new error that satisfies errors.Is(err, ErrEncode).
func Encode(msg string, args ...interface{}) error {
	return newError(ErrEncode, msg, args...)
}

// Existent creates a new error that satisfies errors.Is(err, ErrExistent).
func Existent(msg string, args ...interface{}) error {
	return newError(ErrExistent, msg, args...)
}

// Invalid creates a new error that satisfies errors.Is(err, ErrInvalid).
func Invalid(msg string, args ...interface{}) error {
	return newError(ErrInvalid, msg, args...)
}

// NonExistent creates a new error that satisfies errors.Is(err, ErrNonExistent).
func NonExistent(msg string, args ...interface{}) error {
	return newError(ErrNonExistent, msg, args...)
}

// Timeout creates a new error that satisfies errors.Is(err, ErrTimeout).
func Timeout(msg string, args ...interface{}) error {
	return newError(ErrTimeout, msg, args...)
}

// Transient creates a new error that satisfies errors.Is(err, ErrTransient).
func Transient(msg string, args ...interface{}) error {
	return newError(ErrTransient, msg, args...)
}
