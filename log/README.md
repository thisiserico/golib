# log
--
    import "github.com/thisiserico/golib/log"

Package log is a generated GoMock package.

## Usage

```go
const (
	// PlainFormat uses plain text as output.
	PlainFormat outputFormat = iota

	// JSONFormat uses JSON as output (used by default).
	JSONFormat
)
```

#### type Logger

```go
type Logger interface {
	// Info specifies an informative log entry.
	Info(context.Context, string, Tags)

	// Error specifies an error log entry.
	Error(context.Context, error, Tags)

	// Fatal specifies a fatal log entry.
	Fatal(context.Context, error, Tags)
}
```

Logger defines the used log capabilities.

#### func  NewLogger

```go
func NewLogger(of outputFormat) Logger
```
NewLogger obtains a new logger using the specified configuration.

#### type MockLogger

```go
type MockLogger struct {
}
```

MockLogger is a mock of Logger interface

#### func  NewMockLogger

```go
func NewMockLogger(ctrl *gomock.Controller) *MockLogger
```
NewMockLogger creates a new mock instance

#### func (*MockLogger) EXPECT

```go
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use

#### func (*MockLogger) Error

```go
func (m *MockLogger) Error(arg0 context.Context, arg1 error, arg2 Tags)
```
Error mocks base method

#### func (*MockLogger) Fatal

```go
func (m *MockLogger) Fatal(arg0 context.Context, arg1 error, arg2 Tags)
```
Fatal mocks base method

#### func (*MockLogger) Info

```go
func (m *MockLogger) Info(arg0 context.Context, arg1 string, arg2 Tags)
```
Info mocks base method

#### type MockLoggerMockRecorder

```go
type MockLoggerMockRecorder struct {
}
```

MockLoggerMockRecorder is the mock recorder for MockLogger

#### func (*MockLoggerMockRecorder) Error

```go
func (mr *MockLoggerMockRecorder) Error(arg0, arg1, arg2 interface{}) *gomock.Call
```
Error indicates an expected call of Error

#### func (*MockLoggerMockRecorder) Fatal

```go
func (mr *MockLoggerMockRecorder) Fatal(arg0, arg1, arg2 interface{}) *gomock.Call
```
Fatal indicates an expected call of Fatal

#### func (*MockLoggerMockRecorder) Info

```go
func (mr *MockLoggerMockRecorder) Info(arg0, arg1, arg2 interface{}) *gomock.Call
```
Info indicates an expected call of Info

#### type Tags

```go
type Tags map[constant.Key]constant.Value
```

Tags abstracts a key -> value dictionary.
