# errors
--
    import "github.com/thisiserico/golib/errors"

Package errors provides a way to generate contextual errors.

## Usage

```go
const (
	// Context indicates a context.Context error.
	Context = Category("context-error")

	// Decode indicates that decoding failed.
	Decode = Category("decode")

	// Encode indicates that encoding failed.
	Encode = Category("encode")

	// Existent indicates that an element already exists.
	Existent = Category("existent")

	// Invalid indicates a validation constraint.
	Invalid = Category("invalid")

	// Permanent indicates that the error is permanent – use along Transient.
	Permanent = Category("permanent")

	// NonExistent indicates that an element doesn't exist.
	NonExistent = Category("non-existent")

	// Transient indicates that the error is transient – use along Permanent.
	Transient = Category("transient")
)
```

#### func  Is

```go
func Is(err error, cat Category) bool
```
Is evaluates whether the given error matches a particular category.

#### func  New

```go
func New(args ...interface{}) error
```
New creates a contextual error by accepting different arguments listed below.
The arguments need to be passed in that order to end up with a consistent error.

- `context.Context`

    If the context has erroed, a Context category is added to the list of
    categories.

- `string`

    The given message is added to the error stack.

- `contextualError`

    The existing error stack, categories and tags are preserved and/or extended.

- `error`

    The error message is added to the error stack.

- `Category`

    The given type is stacked. It can be later be accessed with `Is`.

- `kv.Pair`

    To provide more contextual data points.

#### func  Tags

```go
func Tags(err error) []kv.Pair
```
Tags returns the error contextual data points.

#### type Category

```go
type Category string
```

Category indicates the error type on its inception.
