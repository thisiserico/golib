// Package kv provides a way to work with keys and values.
// Those let keep the consistency among packages when working with
// key-value pairs. Functionality for some pre-defined context attributes
// is also provided.
package kv

const redactedValue = "redacted"

// Val encapsulates the given value.
// Keeping it encapsulated allows to work with obfuscated pairs.
// A value can also be used on its own.
type Val struct {
	raw          interface{}
	isObfuscated bool
}

// Pair encapsulates a key-value representation.
// All Val methods can be accessed from a Pair.
type Pair struct {
	key string
	Val
}

// New generates a new Pair using the given key and values.
func New(key string, val interface{}) Pair {
	return Pair{
		key: key,
		Val: Val{raw: val},
	}
}

// NewObfuscated generates a new Pair using the given key. The value, however,
// will be obfuscated. This prevents situations where a value is not supposed
// to be reported to other components. Only strings are supported at this time.
func NewObfuscated(key string, val interface{}) Pair {
	return Pair{
		key: key,
		Val: Val{
			isObfuscated: true,
		},
	}
}

// Name returns the key name of the Pair.
func (p Pair) Name() string {
	return p.key
}

// Value returns a new value to be used outside any Pair.
func Value(v interface{}) Val {
	return Val{raw: v}
}

// Value returns the raw value in its original form. If the value is obfuscated,
// a redacted value is provided instead.
func (v Val) Value() interface{} {
	if v.isObfuscated {
		return redactedValue
	}

	return v.raw
}

// String returns the raw string value, or empty string if one doesn't exist.
// If the value is obfuscated, a redacted value is provided instead.
func (v Val) String() string {
	if v.isObfuscated {
		return redactedValue
	}

	s, ok := v.raw.(string)
	if !ok {
		return ""
	}

	return s
}

// Int returns the raw integer value, or 0 if one doesn't exist.
func (v Val) Int() int {
	i, ok := v.raw.(int)
	if ok {
		return i
	}

	f, ok := v.raw.(float64)
	if ok {
		return int(f)
	}

	return 0
}

// Bool returns the raw boolean value, or false if one doesn't exist.
func (v Val) Bool() bool {
	b, ok := v.raw.(bool)
	if !ok {
		return false
	}

	return b
}
