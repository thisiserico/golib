package logger

import (
	"context"
	"testing"

	"github.com/thisiserico/golib/kv"
	"github.com/thisiserico/golib/logger/memory"
	"github.com/thisiserico/golib/oops"
)

func TestLoggingExecutionAttributes(t *testing.T) {
	const (
		buildID       = "build_id"
		serviceHost   = "service_host"
		serviceName   = "service_name"
		correlationID = "correlation_id"
		isDryRun      = true
	)

	ctx := context.Background()
	ctx = kv.SetStaticAttributes(ctx, buildID, serviceHost, serviceName)
	ctx = kv.SetDynamicAttributes(ctx, correlationID, isDryRun)

	writer := memory.New()
	log := New(writer, JSONOutput)
	log(ctx)

	line, exists := writer.Line(0)
	if !exists {
		t.Fatal("no log line was produced")
	}

	if got := line.Fields["svc.build_id"]; got != buildID {
		t.Fatalf("unexpected build ID, got %s, want %s", got, buildID)
	}
	if got := line.Fields["svc.host"]; got != serviceHost {
		t.Fatalf("unexpected service host, got %s, want %s", got, serviceHost)
	}
	if got := line.Fields["svc.name"]; got != serviceName {
		t.Fatalf("unexpected service name, got %s, want %s", got, serviceName)
	}
	if got := line.Fields["ctx.correlation_id"]; got != correlationID {
		t.Fatalf("unexpected correlation ID, got %s, want %s", got, correlationID)
	}
	if got := line.Fields["ctx.is_dry_run"]; got != true {
		t.Fatalf("unexpected dry run indicator, got %t", got)
	}
}

func TestLoggingAMessage(t *testing.T) {
	const message = "message"

	writer := memory.New()
	log := New(writer, JSONOutput)
	log(message)

	line, exists := writer.Line(0)
	if !exists {
		t.Fatal("no log line was produced")
	}

	if got := line.Message; got != message {
		t.Fatalf("unexpected log line message, got %s, want %s", got, message)
	}
}

func TestLoggingAnError(t *testing.T) {
	t.Run("with only a message", func(t *testing.T) {
		const message = "message"
		err := oops.Invalid(message)

		writer := memory.New()
		log := New(writer, JSONOutput)
		log(err)

		line, exists := writer.Line(0)
		if !exists {
			t.Fatal("no log line was produced")
		}

		if got := line.Message; got != message {
			t.Fatalf("unexpected log line message, got %s, want %s", got, message)
		}
	})

	t.Run("that contains tags", func(t *testing.T) {
		const key = "key"
		const value = "value"
		err := oops.With(oops.Invalid(""), kv.New("key", "value"))

		writer := memory.New()
		log := New(writer, JSONOutput)
		log(err)

		line, exists := writer.Line(0)
		if !exists {
			t.Fatal("no log line was produced")
		}

		if got := line.Fields[key]; got != value {
			t.Fatalf("unexpected log line tag, got %s, want %s", got, value)
		}
	})
}

func TestLoggingTagPairs(t *testing.T) {
	const key = "key"
	const value = "value"

	writer := memory.New()
	log := New(writer, JSONOutput)
	log(kv.New(key, value))

	line, exists := writer.Line(0)
	if !exists {
		t.Fatal("no log line was produced")
	}

	if got := line.Fields[key]; got != value {
		t.Fatalf("unexpected log line tag, got %s, want %s", got, value)
	}
}
