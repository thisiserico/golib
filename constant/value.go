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

// AnyValue composes an unknown value type.
func AnyValue(val interface{}) Value {
	return anyValue{val: val}
}

type anyValue struct {
	val interface{}
}

// Value returns the unknown value.
func (v anyValue) Value() interface{} {
	return v.val
}

// OneBuildID encapsulates a BuildID.
type OneBuildID string

// Value returns the string representation.
func (v OneBuildID) Value() interface{} {
	return string(v)
}

// OneCorrelationID encapsulates a CorrelationID.
type OneCorrelationID string

// Value returns the string representation.
func (v OneCorrelationID) Value() interface{} {
	return string(v)
}

// InBehalfOfServiceName encapsulates an InBehalfOf.
type InBehalfOfServiceName string

// Value returns the string representation.
func (v InBehalfOfServiceName) Value() interface{} {
	return string(v)
}

// IsDryRunExecution encapsulates an IsDryRun.
type IsDryRunExecution bool

// Value returns the boolean representation.
func (v IsDryRunExecution) Value() interface{} {
	return bool(v)
}

// RunningInHost encapsulates a ServiceHost.
type RunningInHost string

// Value returns the string representation.
func (v RunningInHost) Value() interface{} {
	return string(v)
}

// RunningService encapsulates a ServiceName.
type RunningService string

// Value returns the string representation.
func (v RunningService) Value() interface{} {
	return string(v)
}

// RequestedByServiceName encapsulates a WhosRequesting.
type RequestedByServiceName string

// Value returns the string representation.
func (v RequestedByServiceName) Value() interface{} {
	return string(v)
}
