// Package cntxt provides a way to access known context values.
package cntxt

import "context"

type key string

func get(ctx context.Context, k key, zeroValue interface{}) interface{} {
	value := ctx.Value(k)
	if value == nil {
		return zeroValue
	}

	return value
}
