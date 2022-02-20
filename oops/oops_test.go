package oops

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/thisiserico/golib/kv"
)

func TestErrorTypology(t *testing.T) {
	tests := []struct {
		constructor func(string, ...interface{}) error
		typology    error
	}{
		{Cancelled, ErrCancelled},
		{Decode, ErrDecode},
		{Encode, ErrEncode},
		{Existent, ErrExistent},
		{Invalid, ErrInvalid},
		{NonExistent, ErrNonExistent},
		{Timeout, ErrTimeout},
		{Transient, ErrTransient},
	}

	for _, test := range tests {
		t.Run(test.typology.Error(), func(t *testing.T) {
			id := uuid.New().String()
			err := test.constructor("msg: %s", id)

			if !errors.Is(err, test.typology) {
				t.Errorf("the error should be of type %s", test.typology)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		input  error
		output string
	}{
		{
			input:  Cancelled("template message: %w", errors.New("upstream error")),
			output: "upstream error",
		},
		{
			input:  fmt.Errorf("template message: %w", Cancelled("upstream error")),
			output: "upstream error",
		},
		{
			input:  fmt.Errorf("template message: %w", Cancelled("another template: %w", errors.New("upstream error"))),
			output: "another template: upstream error",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			if got := errors.Unwrap(test.input).Error(); got != test.output {
				t.Errorf("unexpected error message, want %s, got %s", test.output, got)
			}
		})
	}
}

func TestContextualDetails(t *testing.T) {
	tests := []struct {
		input   error
		details []kv.Pair
		output  string
	}{
		{
			With(errors.New("oops"), kv.New("key", "value")),
			nil,
			"oops",
		},
		{
			With(Cancelled("oops"), kv.New("first", 1), kv.New("second", "2")),
			[]kv.Pair{
				kv.New("first", 1),
				kv.New("second", "2"),
			},
			"oops",
		},
		{
			With(Cancelled("oops"), kv.New("key", 1), kv.New("key", 2)),
			[]kv.Pair{
				kv.New("key", 1),
				kv.New("key", 2),
			},
			"oops",
		},
		{
			With(With(Cancelled("oops"), kv.New("key", "inner")), kv.New("key", "outer")),
			[]kv.Pair{
				kv.New("key", "inner"),
				kv.New("key", "outer"),
			},
			"oops",
		},
		{
			With(errors.New("oops"), kv.New("key", "value")),
			[]kv.Pair{
				kv.New("key", "value"),
			},
			"oops",
		},
		{
			With(fmt.Errorf("oops: %w", Cancelled("inner")), kv.New("key", "value")),
			[]kv.Pair{
				kv.New("key", "value"),
			},
			"oops: inner",
		},
		{
			With(fmt.Errorf("oops: %w", Cancelled("inner: %w", With(errors.New("inner most"), kv.New("inner", "most")))), kv.New("key", "value")),
			[]kv.Pair{
				kv.New("key", "value"),
				kv.New("inner", "most"),
			},
			"oops: inner: inner most",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			details := Details(test.input)

			for i, pair := range test.details {
				if got := details[i]; got != pair {
					t.Errorf("unexpected pair, want %v, got %v", pair, got)
				}

				if got := test.input.Error(); got != test.output {
					t.Errorf("unexpected error message, want %s, got %s", test.output, got)
				}
			}
		})
	}
}

func TestContextualDetail(t *testing.T) {
	tests := []struct {
		input  error
		key    string
		detail kv.Pair
		exists bool
	}{
		{
			errors.New("oops"),
			"key",
			kv.Pair{},
			false,
		},
		{
			With(errors.New("oops"), kv.New("key", "value")),
			"key",
			kv.New("key", "value"),
			true,
		},
		{
			With(Cancelled("oops"), kv.New("first", 1), kv.New("second", "2")),
			"second",
			kv.New("second", "2"),
			true,
		},
		{
			With(fmt.Errorf("oops: %w", Cancelled("inner: %w", With(errors.New("inner most"), kv.New("inner", "most")))), kv.New("key", "value")),
			"inner",
			kv.New("inner", "most"),
			true,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			pair, exists := Detail(test.input, test.key)

			if want, got := test.detail.String(), pair.String(); got != want {
				t.Errorf("unexpected pair, want %s, got %s", want, got)
			}
			if exists != test.exists {
				t.Errorf("execpected existence indicator, want %t, got %t", test.exists, exists)
			}
		})
	}
}
