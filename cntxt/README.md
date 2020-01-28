# cntxt
--
    import "github.com/thisiserico/golib/cntxt"

Package cntxt provides a way to access known context values.

## Usage

#### func  BuildID

```go
func BuildID(ctx context.Context) kv.Pair
```
BuildID provides a pair with the service build ID stored in the context.

#### func  CorrelationID

```go
func CorrelationID(ctx context.Context) kv.Pair
```
CorrelationID provides a pair with the execution ID stored in the context.

#### func  ExecutionIsDryRun

```go
func ExecutionIsDryRun(ctx context.Context, asDryRun bool) context.Context
```
ExecutionIsDryRun sets the indicator of whether the current execution is a dry
run into the context.

#### func  ExecutionTriggeredBy

```go
func ExecutionTriggeredBy(ctx context.Context, name string) context.Context
```
ExecutionTriggeredBy sets the service name which triggered the current execution
into the context.

#### func  IsDryRun

```go
func IsDryRun(ctx context.Context) kv.Pair
```
IsDryRun provides a pair with the indicator of whether the current execution is
a dry run stored in the context.

#### func  RunningBuildID

```go
func RunningBuildID(ctx context.Context, id string) context.Context
```
RunningBuildID sets the service build ID into the context.

#### func  RunningOnHost

```go
func RunningOnHost(ctx context.Context, host string) context.Context
```
RunningOnHost sets the host where the service is running into the context.

#### func  RunningService

```go
func RunningService(ctx context.Context, service string) context.Context
```
RunningService sets the service name of the service being run into the context.

#### func  ServiceHost

```go
func ServiceHost(ctx context.Context) kv.Pair
```
ServiceHost provides a pair with the service host stored in the context.

#### func  ServiceName

```go
func ServiceName(ctx context.Context) kv.Pair
```
ServiceName provides a pair with the service name stored in the context.

#### func  TriggeredBy

```go
func TriggeredBy(ctx context.Context) kv.Pair
```
TriggeredBy provides a pair with the service name which triggered the execution
stored in the context.

#### func  UsingCorrelationID

```go
func UsingCorrelationID(ctx context.Context, id string) context.Context
```
UsingCorrelationID sets the execution ID into the context.
