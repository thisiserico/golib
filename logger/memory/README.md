# memory
--
    import "github.com/thisiserico/golib/logger/memory"


## Usage

#### type Line

```go
type Line struct {
	Fields  map[string]interface{} `json:"fields"`
	Level   string                 `json:"level"`
	Message string                 `json:"message"`
}
```


#### type Writer

```go
type Writer struct {
}
```


#### func  New

```go
func New() *Writer
```

#### func (*Writer) Line

```go
func (w *Writer) Line(index int) (Line, bool)
```

#### func (*Writer) Write

```go
func (w *Writer) Write(p []byte) (int, error)
```
