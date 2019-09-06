package constant

import "testing"

func Suite(t *testing.T, k Key, v Value, expected interface{}) {
	if value := v.Value(); value != expected {
		t.Fatalf("%s -> want %v, got %v", k, expected, value)
	}
}

func TestConstants(t *testing.T) {
	Suite(t, Key("any value"), AnyValue("any value"), "any value")
	Suite(t, BuildID, OneBuildID("build ID"), "build ID")
	Suite(t, CorrelationID, OneCorrelationID("correlation ID"), "correlation ID")
	Suite(t, InBehalfOf, InBehalfOfServiceName("service name"), "service name")
	Suite(t, IsDryRun, IsDryRunExecution(true), true)
	Suite(t, ServiceHost, RunningInHost("service host"), "service host")
	Suite(t, ServiceName, RunningService("service name"), "service name")
	Suite(t, WhosRequesting, RequestedByServiceName("service name"), "service name")
}
