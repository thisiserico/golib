# constant
--
    import "github.com/thisiserico/golib/constant"

Package constant allows to keep the consistency between data point values by
exposing wide spread keys and value abstractions.

## Usage

```go
const (
	// BuildID indicates the ID for this specific service build.
	BuildID = Key("build_id")

	// CorrelationID holds the unique ID that allows to trace an execution.
	CorrelationID = Key("correlation_id")

	// InBehalfOf indicates the piece of software that launched the execution.
	// Best used along with WhosRequesting.
	InBehalfOf = Key("in_behalf_of")

	// IsDryRun can be used to determine whether a request is a dry run.
	IsDryRun = Key("is_dry_run")

	// ServiceHost indicates the host where the service is currently running.
	ServiceHost = Key("service_host")

	// ServiceName indicates the service name.
	ServiceName = Key("service_name")

	// WhosRequesting indicates the piece of software that is requesting
	// the current execution. Best used along with InBehalfOf.
	WhosRequesting = Key("whos_requesting")
)
```

#### type InBehalfOfServiceName

```go
type InBehalfOfServiceName string
```

InBehalfOfServiceName encapsulates an InBehalfOf.

#### func (InBehalfOfServiceName) Value

```go
func (v InBehalfOfServiceName) Value() interface{}
```
Value returns the string representation.

#### type IsDryRunExecution

```go
type IsDryRunExecution bool
```

IsDryRunExecution encapsulates an IsDryRun.

#### func (IsDryRunExecution) Value

```go
func (v IsDryRunExecution) Value() interface{}
```
Value returns the boolean representation.

#### type Key

```go
type Key string
```

Key encapsulates well known data points used all over the place.

#### type OneBuildID

```go
type OneBuildID string
```

OneBuildID encapsulates a BuildID.

#### func (OneBuildID) Value

```go
func (v OneBuildID) Value() interface{}
```
Value returns the string representation.

#### type OneCorrelationID

```go
type OneCorrelationID string
```

OneCorrelationID encapsulates a CorrelationID.

#### func (OneCorrelationID) Value

```go
func (v OneCorrelationID) Value() interface{}
```
Value returns the string representation.

#### type RequestedByServiceName

```go
type RequestedByServiceName string
```

RequestedByServiceName encapsulates a WhosRequesting.

#### func (RequestedByServiceName) Value

```go
func (v RequestedByServiceName) Value() interface{}
```
Value returns the string representation.

#### type RunningInHost

```go
type RunningInHost string
```

RunningInHost encapsulates a ServiceHost.

#### func (RunningInHost) Value

```go
func (v RunningInHost) Value() interface{}
```
Value returns the string representation.

#### type RunningService

```go
type RunningService string
```

RunningService encapsulates a ServiceName.

#### func (RunningService) Value

```go
func (v RunningService) Value() interface{}
```
Value returns the string representation.

#### type Value

```go
type Value interface {
	// Value returns a typed value.
	Value() interface{}
}
```

Value defines a common way to access constant values.

#### func  AnyValue

```go
func AnyValue(val interface{}) Value
```
AnyValue composes an unknown value type.
