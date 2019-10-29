# halt
--
    import "github.com/thisiserico/golib/halt"

Package halt exposes a convenience method to deal with grafecul shutdowns.

## Usage

#### type Halter

```go
type Halter interface {
	// Wait should block until a shutdown is requested.
	Wait()
}
```

Halter will be used to wait for shutdown requests.

#### func  New

```go
func New(ctx context.Context, logger log.Logger) (context.Context, Halter)
```
New configures and returns the context.Context to use when shutting down.
