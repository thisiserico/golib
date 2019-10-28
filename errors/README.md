# errors
--
    import "github.com/thisiserico/golib/errors"

Package errors provides a way to generate contextual errors.

## Usage

```go
const (
	// Contextual indicates that the error contains contextual tags.
	Contextual = Type(iota)

	// ContextError indicates a context.Context error.
	ContextError

	// Decode indicates that decoding failed.
	Decode

	// Encode indicates that encoding failed.
	Encode

	// Existent indicates that an element already exists.
	Existent

	// Invalid indicates a validatio constraint.
	Invalid

	// Permanent indicates that the error is permanent – use along Transient.
	Permanent

	// PlainError indicates that the error is not contextual.
	PlainError

	// NonExistent indicates that an element doesn't exist.
	NonExistent

	// Transient indicates that the error is transient – use along Permanent.
	Transient
)
```

#### func  Is

```go
func Is(anyError error, requested Type) bool
```
Is returns true when the given error stack contains the requested type.

#### func  New

```go
func New(args ...interface{}) error
```
New facilitates the contextual error creation by accepting different argument
types: context, error, message, type and tags.

- `nil`

    Getting a nil `nil` value explicetely means that there was no error.

- `context.Context`

    If the context has erroed, a ContextError type is added to the list of
    types. The tags will be populated with the known contextual values and
    Contextual type.

- `string`

    The given message is added to the error stack.

- `contextualError`

    To keep the consistency between errors.

- `error`

    The error message is added to the error stack.

- `Type`

    The given type is stacked. It can be later be accessed with `Is`.

- `TagPair`

    To provide more contextual data points.

#### func  Tags

```go
func Tags(anyError error) map[constant.Key]constant.Value
```
Tags returns a key-value dictionary.

#### type TagPair

```go
type TagPair struct {
}
```

TagPair encapsulates a key value pair.

#### func  Tag

```go
func Tag(k constant.Key, v constant.Value) TagPair
```
Tag pairs a key with a value.

#### type Type

```go
type Type int
```

Type indicates the error type on its inception.
