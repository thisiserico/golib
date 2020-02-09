# logger
--
    import "github.com/thisiserico/golib/logger"

Package logger provides a simplified logger to write log lines into a writer.

## Usage

#### type Log

```go
type Log func(...interface{})
```

Log lets clients output log lines into the previously specified writer.
Different argument types –listed below– will be used to define the log line
composition. The order is important, as arguments can override. By default, info
log lines are provided.

- `context.Context`

    Known execution indicators are extracted from the context and provided in
    the log line as tags.

- `string`

    The argument will be used as the log message.

- `error`

    The error message will be used as the log message. An error log line will
    be provided. Error tags will be extracted and used as tags.

- `kv.Pair`

    Each pair will be used as a log line tag.

Other types will be ignored.

#### func  New

```go
func New(w io.Writer, o Output) Log
```
New provides a new logging method. When used, the output will be sent to the
indicated writer, previously formatting the log line using the specified output
method.

#### type Output

```go
type Output int
```

Output defines a way to format log lines before sending them to the writer.

```go
const (
	// PlainOutput writes plain text into the log writer.
	PlainOutput Output = iota

	// JSONOutput writes json encoded text into the log writer.
	JSONOutput
)
```
