// Package pubsubtest simplifies the interaction with the pubsub package
// when running tests.
package pubsubtest

import (
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/thisiserico/golib/pubsub"
)

// IsExpectedEvent will return true when the given event matches with the
// expected name and payload contents.
func IsExpectedEvent(event pubsub.Event, name string, want interface{}) bool {
	if event.ID == "" {
		return false
	}

	if event.Name != pubsub.Name(name) {
		return false
	}

	wantPayload, err := json.Marshal(want)
	if err != nil {
		return false
	}

	var got interface{}
	if err := json.Unmarshal(event.Payload, &got); err != nil {
		return false
	}

	gotPayload, err := json.Marshal(got)
	if err != nil {
		return false
	}

	if diff := cmp.Diff(wantPayload, gotPayload); diff != "" {
		return false
	}

	return true
}
