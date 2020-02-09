package kv

import (
	"errors"
	"testing"
)

func TestPair(t *testing.T) {
	t.Run("using the raw value", func(t *testing.T) {
		key := "key"
		value := "value"
		p := New(key, value)

		if got := p.Name(); key != got {
			t.Fatalf("unexpected pair name, want %s, got %s", key, got)
		}
		if got := p.String(); value != got {
			t.Fatalf("unexpected pair value, want %s, got %s", value, got)
		}
	})

	t.Run("using an obfuscated raw value", func(t *testing.T) {
		key := "key"
		value := errors.New("error")
		p := NewObfuscated(key, value)

		if got := p.Name(); key != got {
			t.Fatalf("unexpected pair name, want %s, got %s", key, got)
		}
		if got := p.Value(); redactedValue != got {
			t.Fatalf("unexpected obfuscated pair value, want %s, got %s", redactedValue, got)
		}
	})

	t.Run("using an obfuscated value", func(t *testing.T) {
		key := "key"
		value := "value"
		p := NewObfuscated(key, value)

		if got := p.Name(); key != got {
			t.Fatalf("unexpected pair name, want %s, got %s", key, got)
		}
		if got := p.String(); redactedValue != got {
			t.Fatalf("unexpected obfuscated pair value, want %s, got %s", redactedValue, got)
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("in its original form", func(t *testing.T) {
		want := errors.New("error")
		v := Value(want)

		if got := v.Value(); want != got {
			t.Fatalf("unexpected raw value, want %s, got %s", want, got)
		}
	})

	t.Run("as string", func(t *testing.T) {
		want := "string"
		v := Value(want)

		if got := v.String(); want != got {
			t.Fatalf("unexpected string, want %s, got %s", want, got)
		}
	})

	t.Run("as int", func(t *testing.T) {
		want := 24
		v := Value(want)

		if got := v.Int(); want != got {
			t.Fatalf("unexpected int, want %d, got %d", want, got)
		}
	})

	t.Run("as bool", func(t *testing.T) {
		want := true
		v := Value(want)

		if got := v.Bool(); want != got {
			t.Fatalf("unexpected bool, want %t, got %t", want, got)
		}
	})
}
