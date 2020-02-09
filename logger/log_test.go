package logger

import (
	"context"
	"testing"

	"github.com/thisiserico/golib/v2/cntxt"
	"github.com/thisiserico/golib/v2/errors"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger/memory"
)

func TestLoggingExecutionAttributes(t *testing.T) {
	const (
		buildID       = "build_id"
		serviceHost   = "service_host"
		serviceName   = "service_name"
		correlationID = "correlation_id"
		triggeredBy   = "triggered_by"
		isDryRun      = true
	)

	ctx := cntxt.RunningBuildID(context.Background(), buildID)
	ctx = cntxt.RunningOnHost(ctx, serviceHost)
	ctx = cntxt.RunningService(ctx, serviceName)
	ctx = cntxt.UsingCorrelationID(ctx, correlationID)
	ctx = cntxt.ExecutionTriggeredBy(ctx, triggeredBy)
	ctx = cntxt.ExecutionIsDryRun(ctx, isDryRun)

	writer := memory.New()
	log := New(writer, JSONOutput)
	log(ctx)

	line, exists := writer.Line(0)
	if !exists {
		t.Fatal("no log line was produced")
	}

	if got := line.Fields["build_id"]; got != buildID {
		t.Fatalf("unexpected build ID, got %s, want %s", got, buildID)
	}
	if got := line.Fields["service_host"]; got != serviceHost {
		t.Fatalf("unexpected service host, got %s, want %s", got, serviceHost)
	}
	if got := line.Fields["service_name"]; got != serviceName {
		t.Fatalf("unexpected service name, got %s, want %s", got, serviceName)
	}
	if got := line.Fields["correlation_id"]; got != correlationID {
		t.Fatalf("unexpected correlation ID, got %s, want %s", got, correlationID)
	}
	if got := line.Fields["triggered_by"]; got != triggeredBy {
		t.Fatalf("unexpected triggered by, got %s, want %s", got, triggeredBy)
	}
	if got := line.Fields["is_dry_run"]; got != true {
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
		err := errors.New(message)

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
		err := errors.New(kv.New("key", "value"))

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
