package kv

const (
	buildID     = "build_id"
	serviceHost = "service_host"
	serviceName = "service_name"

	correlationID = "correlation_id"
	triggeredBy   = "triggered_by"
	isDryRun      = "is_dry_run"
)

// BuildID provides a new Pair that encapsulates a build ID.
func BuildID(id string) Pair {
	return New(buildID, id)
}

// ServiceHost provides a new Pair that encapsulates a service host.
func ServiceHost(host string) Pair {
	return New(serviceHost, host)
}

// ServiceName provides a new Pair that encapsulates a service name.
func ServiceName(name string) Pair {
	return New(serviceName, name)
}

// CorrelationID provides a new Pair that encapsulates a correlation ID.
func CorrelationID(id string) Pair {
	return New(correlationID, id)
}

// TriggeredBy provides a new Pair that encapsulates a triggered by value.
func TriggeredBy(name string) Pair {
	return New(triggeredBy, name)
}

// IsDryRun provides a new Pair that encapsulates an is dry run indicator.
func IsDryRun(dryRun bool) Pair {
	return New(isDryRun, dryRun)
}
