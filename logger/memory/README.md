# memory
--
    import "github.com/thisiserico/golib/logger/memory"

Package memory is a io.Writer implementation to use when testing log lines being
produced.

## Usage

#### type Line

```go
type Line struct {
	// Fields contains the log tags.
	Fields map[string]interface{} `json:"fields"`

	// Level indicates the log level.
	Level string `json:"level"`

	// Message contains the actual message string.
	Message string `json:"message"`
}
```

Line encapsulates the different elements that were logged.

#### type Writer

```go
type Writer struct {
}
```

Writer implements io.Writer and provides a way to fetch the log lines that were
produced.

#### func  New

```go
func New() *Writer
```
New returns a new Writer.

#### func (*Writer) Line

```go
func (w *Writer) Line(index int) (Line, bool)
```
Line fetches the indicated log line. It also returns a boolean indicating
whether the requested log line was produced.

#### func (*Writer) Write

```go
func (w *Writer) Write(p []byte) (int, error)
```
