package dummy

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/thisiserico/golib/v2/cntxt"
	"github.com/thisiserico/golib/v2/kv"
	"github.com/thisiserico/golib/v2/logger"
	"github.com/thisiserico/golib/v2/logger/memory"
)

func TestThatSpansCanBeObtained(t *testing.T) {
	writer := memory.New()
	log := logger.New(writer, logger.JSONOutput)
	ag := Agent(log)

	t.Run("from scratch", func(t *testing.T) {
		spanName := uuid.New().String()
		_, s := ag.StartSpan(context.Background(), spanName)

		if s == nil {
			t.Fatal("a span had to be created")
		}
		if name := s.(*span).name(); name != spanName {
			t.Fatalf("invalid span name, want %s, got %s", spanName, name)
		}
	})

	t.Run("when no span was created yet", func(t *testing.T) {
		const spanName = "new"

		s := ag.GetSpan(context.Background())

		if s == nil {
			t.Fatal("a span had to be created")
		}
		if name := s.(*span).name(); name != spanName {
			t.Fatalf("invalid span name, want %s, got %s", spanName, name)
		}
	})

	t.Run("when a span was already created", func(t *testing.T) {
		spanName := uuid.New().String()
		ctx, _ := ag.StartSpan(context.Background(), spanName)
		s := ag.GetSpan(ctx)

		if s == nil {
			t.Fatal("a span had to be created")
		}
		if name := s.(*span).name(); name != spanName {
			t.Fatalf("invalid span name, want %s, got %s", spanName, name)
		}
	})
}

func TestThatSpansAreReported(t *testing.T) {
	t.Run("once they are completed", func(t *testing.T) {
		writer := memory.New()
		log := logger.New(writer, logger.JSONOutput)
		ag := Agent(log)

		const (
			firstSpan  = "first"
			secondSpan = "second"
		)

		ctx, first := ag.StartSpan(context.Background(), firstSpan)
		ctx, second := ag.StartSpan(ctx, secondSpan)
		second.Complete()
		first.Complete()

		firstLine, exists := writer.Line(0)
		if !exists {
			t.Fatal("a log line should have been written")
		}
		secondLine, exists := writer.Line(1)
		if !exists {
			t.Fatal("a log line should have been written")
		}
		if _, exists := writer.Line(2); exists {
			t.Fatal("only two log lines should have been written")
		}

		expectedFirstLine := fmt.Sprintf("[ %s ]", secondSpan)
		expectedSecondLine := fmt.Sprintf("[ %s ]", firstSpan)
		if got := firstLine.Message; got != expectedFirstLine {
			t.Fatalf("unexpected first line, want %s, got %s", expectedFirstLine, got)
		}
		if got := secondLine.Message; got != expectedSecondLine {
			t.Fatalf("unexpected second line, want %s, got %s", expectedSecondLine, got)
		}
	})

	t.Run("including contextual data", func(t *testing.T) {
		const (
			buildID       = "build_id"
			serviceHost   = "service_host"
			serviceName   = "service_name"
			correlationID = "correlation_id"
			triggeredBy   = "triggered_by"
			isDryRun      = true
		)

		writer := memory.New()
		log := logger.New(writer, logger.JSONOutput)
		ag := Agent(log)

		name := uuid.New().String()
		ctx := cntxt.RunningBuildID(context.Background(), buildID)
		ctx = cntxt.RunningOnHost(ctx, serviceHost)
		ctx = cntxt.RunningService(ctx, serviceName)
		ctx = cntxt.UsingCorrelationID(ctx, correlationID)
		ctx = cntxt.ExecutionTriggeredBy(ctx, triggeredBy)
		ctx = cntxt.ExecutionIsDryRun(ctx, isDryRun)

		_, s := ag.StartSpan(ctx, name)
		s.Complete()

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
	})

	t.Run("including tags", func(t *testing.T) {
		writer := memory.New()
		log := logger.New(writer, logger.JSONOutput)
		ag := Agent(log)

		name := uuid.New().String()
		key := uuid.New().String()
		value := uuid.New().String()

		_, s := ag.StartSpan(context.Background(), name)
		s.AddPair(context.Background(), kv.New(key, value))
		s.Complete()

		line, _ := writer.Line(0)
		val, exists := line.Fields[key]
		if !exists {
			t.Fatal("the field should have been reported")
		}
		if got := val.(string); got != value {
			t.Fatalf("invalid reported field, want %s, got %s", value, got)
		}
	})
}
