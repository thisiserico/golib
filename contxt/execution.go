package contxt

import (
	"context"

	"github.com/thisiserico/golib/constant"
)

// WithCorrelationID adds a correlation ID to the context.
func WithCorrelationID(ctx context.Context, value constant.OneCorrelationID) context.Context {
	return context.WithValue(ctx, constant.CorrelationID, value)
}

// CorrelationID returns the current correlation ID.
func CorrelationID(ctx context.Context) constant.OneCorrelationID {
	val := ctx.Value(constant.CorrelationID)
	if val == nil {
		val = constant.OneCorrelationID("")
	}

	return val.(constant.OneCorrelationID)
}

// WithInBehalfOf informs about the execution launcher to the context.
func WithInBehalfOf(ctx context.Context, value constant.InBehalfOfServiceName) context.Context {
	return context.WithValue(ctx, constant.InBehalfOf, value)
}

// InBehalfOfService returns the current in behalf of service name.
func InBehalfOfService(ctx context.Context) constant.InBehalfOfServiceName {
	val := ctx.Value(constant.InBehalfOf)
	if val == nil {
		val = constant.InBehalfOfServiceName("")
	}

	return val.(constant.InBehalfOfServiceName)
}

// WithWhosRequesting informs about the immediate previous executor to the context.
func WithWhosRequesting(ctx context.Context, value constant.RequestedByServiceName) context.Context {
	return context.WithValue(ctx, constant.WhosRequesting, value)
}

// RequestedByService returns the current who's requesting service name.
func RequestedByService(ctx context.Context) constant.RequestedByServiceName {
	val := ctx.Value(constant.WhosRequesting)
	if val == nil {
		val = constant.RequestedByServiceName("")
	}

	return val.(constant.RequestedByServiceName)
}

// WithIsDryRun informs about the nature of the request to the context.
func WithIsDryRun(ctx context.Context, value constant.IsDryRunExecution) context.Context {
	return context.WithValue(ctx, constant.IsDryRun, value)
}

// IsDryRunExecution returns true when the current execution is a dry run.
func IsDryRunExecution(ctx context.Context) constant.IsDryRunExecution {
	val := ctx.Value(constant.IsDryRun)
	if val == nil {
		val = constant.IsDryRunExecution(false)
	}

	return val.(constant.IsDryRunExecution)
}
