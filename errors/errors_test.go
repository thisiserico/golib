package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/thisiserico/golib/v2/kv"
)

const plainString = "plain string"

var errPlain = errors.New("plain error")

func TestTheContext(t *testing.T) {
	t.Run("when it's canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := New(ctx)
		if !Is(err, Context) {
			t.Fatalf("category doesn't match, want %v, got %v", Context, err)
		}
	})

	t.Run("when it's active", func(t *testing.T) {
		err := New(context.Background())
		if Is(err, Context) {
			t.Fatal("a context category wasn't exppected")
		}
	})
}

func TestTheErrorMessageChain(t *testing.T) {
	t.Run("using a string", func(t *testing.T) {
		err := New(plainString)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}
		if err.Error() != plainString {
			t.Fatalf("unexpected error message, want %s, got %s", plainString, err)
		}
	})

	t.Run("using a plain error", func(t *testing.T) {
		err := New(errPlain)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}
		if err.Error() != errPlain.Error() {
			t.Fatalf("unexpected error message, want %s, got %s", errPlain, err)
		}
	})

	t.Run("using a plain string and an error", func(t *testing.T) {
		err := New(plainString, errPlain)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}
		want := fmt.Sprintf("%s: %s", plainString, errPlain.Error())
		if err.Error() != want {
			t.Fatalf("unexpected error message, want %s, got %s", want, err)
		}
	})

	t.Run("using another contextual error", func(t *testing.T) {
		err := New(plainString, errPlain)
		err = New(plainString, err)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}
		want := fmt.Sprintf("%s: %s: %s", plainString, plainString, errPlain.Error())
		if err.Error() != want {
			t.Fatalf("unexpected error message, want %s, got %s", want, err)
		}
	})
}

func TestCategories(t *testing.T) {
	tests := []struct {
		category Category
	}{
		{category: Context},
		{category: Decode},
		{category: Encode},
		{category: Existent},
		{category: Invalid},
		{category: Permanent},
		{category: NonExistent},
		{category: Transient},
	}

	for _, test := range tests {
		t.Run(string(test.category), func(t *testing.T) {
			err := New(test.category)
			if !Is(err, test.category) {
				t.Fatalf("%v category was expected", test.category)
			}
		})
	}
}

func TestTheTagging(t *testing.T) {
	t.Run("with a single tag", func(t *testing.T) {
		const key = "key"
		const val = "val"

		err := New(kv.New(key, val))
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}
		tags := Tags(err)
		if got := tags[0].Name(); got != key {
			t.Fatalf("invalid tag key, want %s got %s", key, got)
		}
		if got := tags[0].String(); got != val {
			t.Fatalf("invalid tag value, want %s got %s", val, got)
		}
	})
}
