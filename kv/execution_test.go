package kv

import (
	"testing"

	"github.com/lucsky/cuid"
)

func TestExecution(t *testing.T) {
	t.Run("build id", func(t *testing.T) {
		want := cuid.New()
		p := BuildID(want)

		if got := p.String(); want != got {
			t.Fatalf("unexpected build ID, want %s, got %s", want, got)
		}
	})

	t.Run("service host", func(t *testing.T) {
		want := cuid.New()
		p := ServiceHost(want)

		if got := p.String(); want != got {
			t.Fatalf("unexpected service host, want %s, got %s", want, got)
		}
	})

	t.Run("service name", func(t *testing.T) {
		want := cuid.New()
		p := ServiceName(want)

		if got := p.String(); want != got {
			t.Fatalf("unexpected service name, want %s, got %s", want, got)
		}
	})

	t.Run("correlation id", func(t *testing.T) {
		want := cuid.New()
		p := CorrelationID(want)

		if got := p.String(); want != got {
			t.Fatalf("unexpected correlation id, want %s, got %s", want, got)
		}
	})

	t.Run("triggered by", func(t *testing.T) {
		want := cuid.New()
		p := TriggeredBy(want)

		if got := p.String(); want != got {
			t.Fatalf("unexpected triggered by, want %s, got %s", want, got)
		}
	})

	t.Run("is dry run", func(t *testing.T) {
		want := true
		p := IsDryRun(want)

		if got := p.Bool(); want != got {
			t.Fatalf("unexpected is dry run, want %t, got %t", want, got)
		}
	})
}
