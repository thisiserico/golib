# halt
--
    import "github.com/thisiserico/golib/halt"

Package halt exposes a convenience method to deal with graceful shutdowns.

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
func New(ctx context.Context, log logger.Log) (context.Context, Halter)
```
New configures and returns the context to use when shutting down.
