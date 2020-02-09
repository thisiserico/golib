package cntxt

import (
	"context"

	"github.com/thisiserico/golib/v2/kv"
)

var (
	buildID     = key("build_id")
	serviceHost = key("service_host")
	serviceName = key("service_name")

	correlationID = key("correlation_id")
	triggeredBy   = key("triggered_by")
	isDryRun      = key("is_dry_run")
)

// RunningBuildID sets the service build ID into the context.
func RunningBuildID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, buildID, kv.BuildID(id))
}

// RunningOnHost sets the host where the service is running into the context.
func RunningOnHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, serviceHost, kv.ServiceHost(host))
}

// RunningService sets the service name of the service being run into the context.
func RunningService(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceName, kv.ServiceName(service))
}

// UsingCorrelationID sets the execution ID into the context.
func UsingCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationID, kv.CorrelationID(id))
}

// ExecutionTriggeredBy sets the service name which triggered the current
// execution into the context.
func ExecutionTriggeredBy(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, triggeredBy, kv.TriggeredBy(name))
}

// ExecutionIsDryRun sets the indicator of whether the current execution is
// a dry run into the context.
func ExecutionIsDryRun(ctx context.Context, asDryRun bool) context.Context {
	return context.WithValue(ctx, isDryRun, kv.IsDryRun(asDryRun))
}

// BuildID provides a pair with the service build ID stored in the context.
func BuildID(ctx context.Context) kv.Pair {
	pair := get(ctx, buildID, nil)
	if pair == nil {
		return kv.BuildID("")
	}

	return pair.(kv.Pair)
}

// ServiceHost provides a pair with the service host stored in the context.
func ServiceHost(ctx context.Context) kv.Pair {
	pair := get(ctx, serviceHost, nil)
	if pair == nil {
		return kv.ServiceHost("")
	}

	return pair.(kv.Pair)
}

// ServiceName provides a pair with the service name stored in the context.
func ServiceName(ctx context.Context) kv.Pair {
	pair := get(ctx, serviceName, nil)
	if pair == nil {
		return kv.ServiceName("")
	}

	return pair.(kv.Pair)
}

// CorrelationID provides a pair with the execution ID stored in the context.
func CorrelationID(ctx context.Context) kv.Pair {
	pair := get(ctx, correlationID, nil)
	if pair == nil {
		return kv.CorrelationID("")
	}

	return pair.(kv.Pair)
}

// TriggeredBy provides a pair with the service name which triggered the
// execution stored in the context.
func TriggeredBy(ctx context.Context) kv.Pair {
	pair := get(ctx, triggeredBy, nil)
	if pair == nil {
		return kv.TriggeredBy("")
	}

	return pair.(kv.Pair)
}

// IsDryRun provides a pair with the indicator of whether the current
// execution is a dry run stored in the context.
func IsDryRun(ctx context.Context) kv.Pair {
	pair := get(ctx, isDryRun, nil)
	if pair == nil {
		return kv.IsDryRun(false)
	}

	return pair.(kv.Pair)
}
