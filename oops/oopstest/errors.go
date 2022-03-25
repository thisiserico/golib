package oopstest

import "errors"

// Is can be used like errors.Is, but matching several targets at once.
// The error needs to match all targets for this to evaluate to true.
func Is(err error, targets ...error) bool {
	for _, target := range targets {
		if !errors.Is(err, target) {
			return false
		}
	}

	return true
}
