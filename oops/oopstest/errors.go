package oopstest

import "errors"

func Is(err error, targets ...error) bool {
	for _, target := range targets {
		if !errors.Is(err, target) {
			return false
		}
	}

	return true
}
