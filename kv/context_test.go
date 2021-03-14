package kv

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestContextPairs(t *testing.T) {
	t.Run("setting all attributes", func(t *testing.T) {
		var (
			buildID       = uuid.New().String()
			serviceHost   = uuid.New().String()
			serviceName   = uuid.New().String()
			correlationID = uuid.New().String()
			isDryRun      = true
		)

		ctx := context.Background()
		ctx = SetStaticAttributes(ctx, buildID, serviceHost, serviceName)
		ctx = SetDynamicAttributes(ctx, correlationID, isDryRun)

		for _, attr := range AllAttributes(ctx) {
			var want interface{}
			var got interface{}

			switch attr.Name() {
			case "build_id":
				want = buildID
				got = attr.String()
			case "service_host":
				want = serviceHost
				got = attr.String()
			case "service_name":
				want = serviceName
				got = attr.String()
			case "correlation_id":
				want = correlationID
				got = attr.String()
			case "is_dry_run":
				want = isDryRun
				got = attr.Bool()
			}

			if want != got {
				t.Fatalf("unexpected attribute, want %s, got %s", want, got)
			}
		}
	})

	t.Run("decorating context", func(t *testing.T) {
		var (
			buildID     = uuid.New().String()
			serviceHost = uuid.New().String()
			serviceName = uuid.New().String()
		)

		ctx := context.Background()
		ctx = SetStaticAttributes(ctx, buildID, serviceHost, serviceName)
		ctx = DecorateWithAttributes(context.Background(), ctx)

		for _, attr := range AllAttributes(ctx) {
			var want interface{}
			var got interface{}

			switch attr.Name() {
			case "build_id":
				want = buildID
				got = attr.String()
			case "service_host":
				want = serviceHost
				got = attr.String()
			case "service_name":
				want = serviceName
				got = attr.String()
			}

			if want != got {
				t.Fatalf("unexpected attribute, want %s, got %s", want, got)
			}
		}
	})

	t.Run("setting dynamic attributes", func(t *testing.T) {
		var (
			correlationID = uuid.New().String()
			isDryRun      = true
		)

		ctx := SetDynamicAttributes(context.Background(), correlationID, isDryRun)

		if want, got := correlationID, CorrelationID(ctx).String(); want != got {
			t.Fatalf("unexpected correlation id, want %s, got %s", want, got)
		}
		if want, got := isDryRun, IsDryRun(ctx).Bool(); want != got {
			t.Fatalf("unexpected is dry run, want %t, got %t", want, got)
		}
	})
}
