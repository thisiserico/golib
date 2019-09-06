package contxt

import (
	"context"

	"github.com/thisiserico/golib/constant"
)

// RunningOnBuildID specifies the build ID being run to the context.
func RunningOnBuildID(ctx context.Context, value constant.OneBuildID) context.Context {
	return context.WithValue(ctx, constant.BuildID, value)
}

// BuildID returns the build ID being run.
func BuildID(ctx context.Context) constant.OneBuildID {
	val := ctx.Value(constant.BuildID)
	if val == nil {
		val = constant.OneBuildID("")
	}

	return val.(constant.OneBuildID)
}

// HostBeingUsed specifies the host being used to the context.
func HostBeingUsed(ctx context.Context, value constant.RunningInHost) context.Context {
	return context.WithValue(ctx, constant.ServiceHost, value)
}

// RunningInHost returns the host being used.
func RunningInHost(ctx context.Context) constant.RunningInHost {
	val := ctx.Value(constant.ServiceHost)
	if val == nil {
		val = constant.RunningInHost("")
	}

	return val.(constant.RunningInHost)
}

// RunningServiceName specifies the service being run to the context.
func RunningServiceName(ctx context.Context, value constant.RunningService) context.Context {
	return context.WithValue(ctx, constant.ServiceName, value)
}

// RunningService returns the service being run.
func RunningService(ctx context.Context) constant.RunningService {
	val := ctx.Value(constant.ServiceName)
	if val == nil {
		val = constant.RunningService("")
	}

	return val.(constant.RunningService)
}
