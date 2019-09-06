package constant

var (
	_ Value = OneBuildID("")
	_ Value = OneCorrelationID("")
	_ Value = InBehalfOfServiceName("")
	_ Value = IsDryRunExecution(false)
	_ Value = RunningInHost("")
	_ Value = RunningService("")
	_ Value = RequestedByServiceName("")
)

// Value defines a common way to access constant values.
type Value interface {
	// Value returns a typed value.
	Value() interface{}
}

// OneBuildID encapsulates a BuildID.
type OneBuildID string

// Value returns the string representation.
func (t OneBuildID) Value() interface{} {
	return string(t)
}

// OneCorrelationID encapsulates a CorrelationID.
type OneCorrelationID string

// Value returns the string representation.
func (t OneCorrelationID) Value() interface{} {
	return string(t)
}

// InBehalfOfServiceName encapsulates an InBehalfOf.
type InBehalfOfServiceName string

// Value returns the string representation.
func (t InBehalfOfServiceName) Value() interface{} {
	return string(t)
}

// IsDryRunExecution encapsulates an IsDryRun.
type IsDryRunExecution bool

// Value returns the boolean representation.
func (t IsDryRunExecution) Value() interface{} {
	return bool(t)
}

// RunningInHost encapsulates a ServiceHost.
type RunningInHost string

// Value returns the string representation.
func (t RunningInHost) Value() interface{} {
	return string(t)
}

// RunningService encapsulates a ServiceName.
type RunningService string

// Value returns the string representation.
func (t RunningService) Value() interface{} {
	return string(t)
}

// RequestedByServiceName encapsulates a WhosRequesting.
type RequestedByServiceName string

// Value returns the string representation.
func (t RequestedByServiceName) Value() interface{} {
	return string(t)
}
