package kv

import "context"

const (
	buildIDKey     = key("build_id")
	serviceHostKey = key("service_host")
	serviceNameKey = key("service_name")

	correlationIDKey = key("correlation_id")
	isDryRunKey      = key("is_dry_run")
)

type key string

func get(ctx context.Context, k key) Pair {
	val := ctx.Value(k)
	if val == nil {
		val = New(string(k), nil)
	}

	return val.(Pair)
}

// DecorateWithAttributes adds the static attributes as values in the
// resulting context.
func DecorateWithAttributes(inUse, background context.Context) context.Context {
	inUse = context.WithValue(inUse, buildIDKey, get(background, buildIDKey))
	inUse = context.WithValue(inUse, serviceHostKey, get(background, serviceHostKey))
	inUse = context.WithValue(inUse, serviceNameKey, get(background, serviceNameKey))

	return inUse
}

// SetStaticAttributes sets program attributes (build ID, service host and service name)
// into the given context.
func SetStaticAttributes(
	ctx context.Context,
	buildID, serviceHost, serviceName string,
) context.Context {
	ctx = context.WithValue(ctx, buildIDKey, New(string(buildIDKey), buildID))
	ctx = context.WithValue(ctx, serviceHostKey, New(string(serviceHostKey), serviceHost))
	ctx = context.WithValue(ctx, serviceNameKey, New(string(serviceNameKey), serviceName))

	return ctx
}

// SetDynamicAttributes sets request attributes (correlation ID and the is dry run)
// into the given context.
func SetDynamicAttributes(
	ctx context.Context,
	correlationID string,
	isDryRun bool,
) context.Context {
	ctx = context.WithValue(ctx, correlationIDKey, New(string(correlationIDKey), correlationID))
	ctx = context.WithValue(ctx, isDryRunKey, New(string(isDryRunKey), isDryRun))

	return ctx
}

// AllAttributes returns all the known pairs that exist in the context.
func AllAttributes(ctx context.Context) []Pair {
	return []Pair{
		get(ctx, buildIDKey),
		get(ctx, serviceHostKey),
		get(ctx, serviceNameKey),
		get(ctx, correlationIDKey),
		get(ctx, isDryRunKey),
	}
}

// CorrelationID returns pair holding that information from the given context.
func CorrelationID(ctx context.Context) Pair {
	return get(ctx, correlationIDKey)
}

// IsDryRun returns pair holding that information from the given context.
func IsDryRun(ctx context.Context) Pair {
	return get(ctx, isDryRunKey)
}
