package cntxt

import (
	"context"
	"testing"
)

func TestRunningBuildIDildID(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if got := BuildID(context.TODO()).String(); got != "" {
			t.Fatalf("no build ID was expected, got %s", got)
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		const want = "build_id"

		ctx := RunningBuildID(context.TODO(), want)
		if got := BuildID(ctx).String(); got != want {
			t.Fatalf("unexpected build ID, want %s, got %s", want, got)
		}
	})
}

func TestServiceHost(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if got := ServiceHost(context.TODO()).String(); got != "" {
			t.Fatalf("no service host was expected, got %s", got)
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		const want = "service_host"

		ctx := RunningOnHost(context.TODO(), want)
		if got := ServiceHost(ctx).String(); got != want {
			t.Fatalf("unexpected service host, want %s, got %s", want, got)
		}
	})
}

func TestServiceName(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if got := ServiceName(context.TODO()).String(); got != "" {
			t.Fatalf("no service name was expected, got %s", got)
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		const want = "service_name"

		ctx := RunningService(context.TODO(), want)
		if got := ServiceName(ctx).String(); got != want {
			t.Fatalf("unexpected service name, want %s, got %s", want, got)
		}
	})
}

func TestCorrelationID(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if got := CorrelationID(context.TODO()).String(); got != "" {
			t.Fatalf("no correlation ID was expected, got %s", got)
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		const want = "correlation_id"

		ctx := UsingCorrelationID(context.TODO(), want)
		if got := CorrelationID(ctx).String(); got != want {
			t.Fatalf("unexpected correlation ID, want %s, got %s", want, got)
		}
	})
}

func TestTriggeredBy(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if got := TriggeredBy(context.TODO()).String(); got != "" {
			t.Fatalf("no triggered by was expected, got %s", got)
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		const want = "triggered_by"

		ctx := ExecutionTriggeredBy(context.TODO(), want)
		if got := TriggeredBy(ctx).String(); got != want {
			t.Fatalf("unexpected triggered by, want %s, got %s", want, got)
		}
	})
}

func TestIsDryRun(t *testing.T) {
	t.Run("when not set in the context", func(t *testing.T) {
		if IsDryRun(context.TODO()).Bool() {
			t.Fatal("no dry run was expected")
		}
	})

	t.Run("when set in the context", func(t *testing.T) {
		ctx := ExecutionIsDryRun(context.TODO(), true)
		if !IsDryRun(ctx).Bool() {
			t.Fatal("a dry run execution was expected")
		}
	})
}
