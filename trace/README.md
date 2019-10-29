# trace
--
    import "github.com/thisiserico/golib/trace"

Package trace let clients handle tracing segments.

## Usage

#### type Segment

```go
type Segment struct {
}
```

Segment encapsulates a tracing span.

#### func  NewSegment

```go
func NewSegment(ctx context.Context, name string) *Segment
```
NewSegment initializes a new tracing segment.

#### func (*Segment) Finish

```go
func (s *Segment) Finish(err *error)
```
Finish finalizes the segment span.

#### func (*Segment) Log

```go
func (s *Segment) Log(key constant.Key, value constant.Value)
```
Log records different events.
