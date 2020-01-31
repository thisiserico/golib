# kv
--
    import "github.com/thisiserico/golib/kv"

Package kv provides a way to work with keys and values. Those let keep the
consistency among packages when working with key-value pairs.

## Usage

#### type Pair

```go
type Pair struct {
	Val
}
```

Pair encapsulates a key-value representation. All Val methods can be accessed
from a Pair.

#### func  BuildID

```go
func BuildID(id string) Pair
```
BuildID provides a new Pair that encapsulates a build ID.

#### func  CorrelationID

```go
func CorrelationID(id string) Pair
```
CorrelationID provides a new Pair that encapsulates a correlation ID.

#### func  IsDryRun

```go
func IsDryRun(dryRun bool) Pair
```
IsDryRun provides a new Pair that encapsulates an is dry run indicator.

#### func  New

```go
func New(key string, val interface{}) Pair
```
New generates a new Pair using the given key and values.

#### func  NewObfuscated

```go
func NewObfuscated(key string, val interface{}) Pair
```
NewObfuscated generates a new Pair using the given key. The value, however, will
be obfuscated. This prevents situations where a value is not supposed to be
reported to other components. Only strings are supported at this time.

#### func  ServiceHost

```go
func ServiceHost(host string) Pair
```
ServiceHost provides a new Pair that encapsulates a service host.

#### func  ServiceName

```go
func ServiceName(name string) Pair
```
ServiceName provides a new Pair that encapsulates a service name.

#### func  TriggeredBy

```go
func TriggeredBy(name string) Pair
```
TriggeredBy provides a new Pair that encapsulates a triggered by value.

#### func (Pair) Name

```go
func (p Pair) Name() string
```
Name returns the key name of the Pair.

#### type Val

```go
type Val struct {
}
```

Val encapsulates the given value. Keeping it encapsulated allows to work with
obfuscated pairs. A value can also be used on its own.

#### func  Value

```go
func Value(v interface{}) Val
```
Value returns the value of the Pair.

#### func (Val) Bool

```go
func (v Val) Bool() bool
```
Bool returns the raw boolean value.

#### func (Val) Int

```go
func (v Val) Int() int
```
Int returns the raw integer value.

#### func (Val) String

```go
func (v Val) String() string
```
String returns the raw string value. If the value is obfuscated, a redacted
value is provided instead.

#### func (Val) Value

```go
func (v Val) Value() interface{}
```
Value returns the raw value in its original form. If the value is obfuscated, a
redacted value is provided instead.
