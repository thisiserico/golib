package contxt

import (
	"context"
	"testing"

	"github.com/thisiserico/golib/constant"
)

func TestExecution(t *testing.T) {
	t.Run("without correlation ID", func(t *testing.T) {
		expected := ""

		if value := string(CorrelationID(context.Background())); value != expected {
			t.Fatalf("invalid correlation ID: want %v, got %v", expected, value)
		}
	})

	t.Run("with correlation ID", func(t *testing.T) {
		expected := "correlation id"
		ctx := WithCorrelationID(context.Background(), constant.OneCorrelationID(expected))

		if value := string(CorrelationID(ctx)); value != expected {
			t.Fatalf("invalid correlation ID: want %v, got %v", expected, value)
		}
	})

	t.Run("without in behalf of", func(t *testing.T) {
		expected := ""

		if value := string(InBehalfOfService(context.Background())); value != expected {
			t.Fatalf("invalid in behalf of: want %v, got %v", expected, value)
		}
	})

	t.Run("with in behalf of", func(t *testing.T) {
		expected := "in behalf of"
		ctx := WithInBehalfOf(context.Background(), constant.InBehalfOfServiceName(expected))

		if value := string(InBehalfOfService(ctx)); value != expected {
			t.Fatalf("invalid in behalf of: want %v, got %v", expected, value)
		}
	})

	t.Run("without who's requesting", func(t *testing.T) {
		expected := ""

		if value := string(RequestedByService(context.Background())); value != expected {
			t.Fatalf("invalid who's requesting: want %v, got %v", expected, value)
		}
	})

	t.Run("with who's requesting", func(t *testing.T) {
		expected := "who's requesting"
		ctx := WithWhosRequesting(context.Background(), constant.RequestedByServiceName(expected))

		if value := string(RequestedByService(ctx)); value != expected {
			t.Fatalf("invalid who's requesting: want %v, got %v", expected, value)
		}
	})

	t.Run("without is dry run", func(t *testing.T) {
		expected := false

		if value := bool(IsDryRunExecution(context.Background())); value != expected {
			t.Fatalf("invalid is dry run: want %v, got %v", expected, value)
		}
	})

	t.Run("with is dry run", func(t *testing.T) {
		expected := true
		ctx := WithIsDryRun(context.Background(), constant.IsDryRunExecution(expected))

		if value := bool(IsDryRunExecution(ctx)); value != expected {
			t.Fatalf("invalid is dry run: want %v, got %v", expected, value)
		}
	})
}

func TestProcess(t *testing.T) {
	t.Run("without build ID", func(t *testing.T) {
		expected := ""

		if value := string(BuildID(context.Background())); value != expected {
			t.Fatalf("invalid build ID: want %v, got %v", expected, value)
		}
	})

	t.Run("with build ID", func(t *testing.T) {
		expected := "build id"
		ctx := RunningOnBuildID(context.Background(), constant.OneBuildID(expected))

		if value := string(BuildID(ctx)); value != expected {
			t.Fatalf("invalid build ID: want %v, got %v", expected, value)
		}
	})

	t.Run("without service host", func(t *testing.T) {
		expected := ""

		if value := string(RunningInHost(context.Background())); value != expected {
			t.Fatalf("invalid service host: want %v, got %v", expected, value)
		}
	})

	t.Run("with service host", func(t *testing.T) {
		expected := "service host"
		ctx := HostBeingUsed(context.Background(), constant.RunningInHost(expected))

		if value := string(RunningInHost(ctx)); value != expected {
			t.Fatalf("invalid service host: want %v, got %v", expected, value)
		}
	})

	t.Run("without service name", func(t *testing.T) {
		expected := ""

		if value := string(RunningService(context.Background())); value != expected {
			t.Fatalf("invalid service name: want %v, got %v", expected, value)
		}
	})

	t.Run("with service name", func(t *testing.T) {
		expected := "service name"
		ctx := RunningServiceName(context.Background(), constant.RunningService(expected))

		if value := string(RunningService(ctx)); value != expected {
			t.Fatalf("invalid service name: want %v, got %v", expected, value)
		}
	})
}
