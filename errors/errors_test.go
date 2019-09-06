package errors

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/thisiserico/golib/constant"
	"github.com/thisiserico/golib/contxt"
)

const plainString = "plain string"

var errPlain = errors.New("plain error")

func TestWithNil(t *testing.T) {
	if err := New(nil); err != nil {
		t.Fatalf("no errors were expected, got %v", err)
	}
}

func TestTheContext(t *testing.T) {
	t.Run("when it's canceled", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := New(ctx)
		if !Is(err, ContextError) {
			t.Fatalf("type doesn't match, want %v, got %v", ContextError, err)
		}
	})

	t.Run("when it's empty", func(t *testing.T) {
		t.Parallel()

		err := New(context.Background())
		if !Is(err, PlainError) {
			t.Fatalf("type doesn't match, want %v, got %v", PlainError, err)
		}
		if !Is(err, Contextual) {
			t.Fatalf("type doesn't match, want %v, got %v", Contextual, err)
		}
	})

	t.Run("when it's filled up", func(t *testing.T) {
		t.Parallel()

		const (
			buildID            = "build_id"
			host               = "host"
			runningService     = "service"
			correlationID      = "correlation_id"
			initializerService = "initializer_service"
			requestingService  = "requesting_service"
			isDryRun           = true
		)

		ctx := context.Background()
		ctx = contxt.RunningOnBuildID(ctx, constant.OneBuildID(buildID))
		ctx = contxt.HostBeingUsed(ctx, constant.RunningInHost(host))
		ctx = contxt.RunningServiceName(ctx, constant.RunningService(runningService))
		ctx = contxt.WithCorrelationID(ctx, constant.OneCorrelationID(correlationID))
		ctx = contxt.WithInBehalfOf(ctx, constant.InBehalfOfServiceName(initializerService))
		ctx = contxt.WithWhosRequesting(ctx, constant.RequestedByServiceName(requestingService))
		ctx = contxt.WithIsDryRun(ctx, constant.IsDryRunExecution(isDryRun))

		err := New(ctx)
		if !Is(err, Contextual) {
			t.Fatal("a contextual type was expected")
		}

		tags := Tags(err)
		if len(tags) != 7 {
			t.Fatalf("unexpected number of tags, want 7, got %d", len(tags))
		}

		assert := func(key constant.Key, expected interface{}) {
			value, wasSet := tags[key]
			if !wasSet {
				t.Fatalf("the %v wasn't set", key)
			}
			if rawValue := value.Value(); rawValue != expected {
				t.Fatalf("unexpected %v set, want %s, got %s", key, expected, rawValue)
			}
		}

		assert(constant.BuildID, buildID)
		assert(constant.ServiceHost, host)
		assert(constant.ServiceName, runningService)
		assert(constant.CorrelationID, correlationID)
		assert(constant.InBehalfOf, initializerService)
		assert(constant.WhosRequesting, requestingService)
		assert(constant.IsDryRun, isDryRun)
	})
}

func TestTheErrorMessageChain(t *testing.T) {
	t.Run("using a string", func(t *testing.T) {
		t.Parallel()

		err := New(plainString)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}

		if err.Error() != plainString {
			t.Fatalf("unexpected error message, want %s, got %s", plainString, err)
		}
	})

	t.Run("using a plain error", func(t *testing.T) {
		t.Parallel()

		err := New(errPlain)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}

		if err.Error() != errPlain.Error() {
			t.Fatalf("unexpected error message, want %s, got %s", errPlain, err)
		}
	})

	t.Run("using a plain error and a string", func(t *testing.T) {
		t.Parallel()

		err := New(plainString, errPlain)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}

		expectedString := fmt.Sprintf("%s: %s", plainString, errPlain.Error())
		if err.Error() != expectedString {
			t.Fatalf("unexpected error message, want %s, got %s", expectedString, err)
		}
	})

	t.Run("using another contextual error", func(t *testing.T) {
		t.Parallel()

		err := New(plainString, errPlain)
		err = New(plainString, err)
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}

		expectedString := fmt.Sprintf("%s: %s: %s", plainString, plainString, errPlain.Error())
		if err.Error() != expectedString {
			t.Fatalf("unexpected error message, want %s, got %s", expectedString, err)
		}
	})
}

func TestTheTyping(t *testing.T) {
	t.Run("making it transient", func(t *testing.T) {
		t.Parallel()

		err := New(Transient)
		if !Is(err, Transient) {
			t.Fatal("a transient error was expected")
		}
	})

	t.Run("making it permanent", func(t *testing.T) {
		t.Parallel()

		err := New(Permanent)
		if !Is(err, Permanent) {
			t.Fatal("a permanent error was expected")
		}
	})

	t.Run("stacking types", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := New(ctx, Permanent)
		err = New(err, Transient)

		if !Is(err, ContextError) {
			t.Fatal("a context error was expected")
		}
		if !Is(err, Contextual) {
			t.Fatal("a contextual error was expected")
		}
		if !Is(err, Permanent) {
			t.Fatal("a permanent error was expected")
		}
		if !Is(err, Transient) {
			t.Fatal("a transient error was expected")
		}
	})
}

type value string

func (val value) Value() interface{} {
	return string(val)
}

func TestTheTagging(t *testing.T) {
	t.Run("with a single tag", func(t *testing.T) {
		t.Parallel()

		err := New(Tag(constant.Key("key"), value("value")))
		if err == nil {
			t.Fatal("an error was expected, got nil")
		}

		tags := Tags(err)
		if tags[constant.Key("key")] != value("value") {
			t.Fatalf(
				"key %s doesn't match with value, want %v, got %v",
				constant.Key("key"),
				value("value"),
				tags[constant.Key("key")],
			)
		}
	})

	t.Run("overriding tags", func(t *testing.T) {
		t.Parallel()

		err := New(Tag(constant.Key("key"), value("first value")))
		err = New(err, Tag(constant.Key("key"), value("second value")))

		tags := Tags(err)
		if tags[constant.Key("key")] != value("first value") {
			t.Fatalf(
				"key %s doesn't match with value, want %v, got %v",
				constant.Key("key"),
				value("first value"),
				tags[constant.Key("key")],
			)
		}
	})

	t.Run("stacking tags", func(t *testing.T) {
		t.Parallel()

		err := New(Tag(constant.Key("key"), value("first value")))
		err = New(err, Tag(constant.Key("key"), value("second value")))
		err = New(err, Tag(constant.Key("random attribute"), value("random value")))

		tags := Tags(err)
		if tags[constant.Key("key")] != value("first value") {
			t.Fatalf(
				"key %s doesn't match with value, want %v, got %v",
				constant.Key("key"),
				value("first value"),
				tags[constant.Key("key")],
			)
		}
		if tags[constant.Key("random attribute")] != value("random value") {
			t.Fatalf(
				"key %s doesn't match with value, want %v, got %v",
				constant.Key("random attribute"),
				value("random value"),
				tags[constant.Key("random attribute")],
			)
		}
	})
}
