package oops

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
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
