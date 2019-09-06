# contxt
--
    import "github.com/thisiserico/golib/contxt"

Package contxt allows to manipulate a context with generic concerns.

## Usage

#### func  BuildID

```go
func BuildID(ctx context.Context) constant.OneBuildID
```
BuildID returns the build ID being run.

#### func  CorrelationID

```go
func CorrelationID(ctx context.Context) constant.OneCorrelationID
```
CorrelationID returns the current correlation ID.

#### func  HostBeingUsed

```go
func HostBeingUsed(ctx context.Context, value constant.RunningInHost) context.Context
```
HostBeingUsed specifies the host being used to the context.

#### func  InBehalfOfService

```go
func InBehalfOfService(ctx context.Context) constant.InBehalfOfServiceName
```
InBehalfOfService returns the current in behalf of service name.

#### func  IsDryRunExecution

```go
func IsDryRunExecution(ctx context.Context) constant.IsDryRunExecution
```
IsDryRunExecution returns true when the current execution is a dry run.

#### func  RequestedByService

```go
func RequestedByService(ctx context.Context) constant.RequestedByServiceName
```
RequestedByService returns the current who's requesting service name.

#### func  RunningInHost

```go
func RunningInHost(ctx context.Context) constant.RunningInHost
```
RunningInHost returns the host being used.

#### func  RunningOnBuildID

```go
func RunningOnBuildID(ctx context.Context, value constant.OneBuildID) context.Context
```
RunningOnBuildID specifies the build ID being run to the context.

#### func  RunningService

```go
func RunningService(ctx context.Context) constant.RunningService
```
RunningService returns the service being run.

#### func  RunningServiceName

```go
func RunningServiceName(ctx context.Context, value constant.RunningService) context.Context
```
RunningServiceName specifies the service being run to the context.

#### func  WithCorrelationID

```go
func WithCorrelationID(ctx context.Context, value constant.OneCorrelationID) context.Context
```
WithCorrelationID adds a correlation ID to the context.

#### func  WithInBehalfOf

```go
func WithInBehalfOf(ctx context.Context, value constant.InBehalfOfServiceName) context.Context
```
WithInBehalfOf informs about the execution launcher to the context.

#### func  WithIsDryRun

```go
func WithIsDryRun(ctx context.Context, value constant.IsDryRunExecution) context.Context
```
WithIsDryRun informs about the nature of the request to the context.

#### func  WithWhosRequesting

```go
func WithWhosRequesting(ctx context.Context, value constant.RequestedByServiceName) context.Context
```
WithWhosRequesting informs about the immediate previous executor to the context.
