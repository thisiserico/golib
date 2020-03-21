# o11y
--
    import "github.com/thisiserico/golib/o11y"

Package o11y provides an abstraction to make system observable. Although the
unit of work is the span, other nuances are being dealt with underneath, like
logging when necessary.

## Usage

#### func  Register

```go
func Register(a Agent)
```
Register allows clients to specify the agent they want to use under the hood.
All operations will be delegated to a specific agent. Calling this method more
than once will overwrite the previously specified strategy.

#### type Agent

```go
type Agent interface {
	// StartSpan needs to be used when a new portion of work just started.
	StartSpan(context.Context, string) (context.Context, Span)

	// GetSpan on the other hand, is encouraged to be used when the program
	// is still dealing with the same unit of work.
	GetSpan(context.Context) Span
}
```

Agent defines a way to either materialize an existing span in the context or to
create a new one from scratch.

#### type Span

```go
type Span interface {
	// AddField allows to provide such context.
	AddField(context.Context, kv.Pair)

	// Send is used when the unit of work has completed.
	Send()
}
```

Span defines a way to inform the system about the context of a certain
execution.

#### func  GetSpan

```go
func GetSpan(ctx context.Context) Span
```
GetSpan on the other hand, tries to extract an existing span from the given
context. Each agent will deal with the fact that a previous span might not exist
yet.

#### func  StartSpan

```go
func StartSpan(ctx context.Context, name string) (context.Context, Span)
```
StartSpan generates a new span using the specified agent. Each agent will have
to provide the span in the context.
