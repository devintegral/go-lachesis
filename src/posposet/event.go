package posposet

import (
	"fmt"
)

// Event is a poset event.
type Event struct {
	Creator Address
	Parents EventHashes

	hash    EventHash
	parents map[EventHash]*Event
}

// Hash calcs hash of event.
func (e *Event) Hash() EventHash {
	if e.hash.IsZero() {
		e.hash = EventHashOf(e)
	}
	return e.hash
}

// String returns string representation.
func (e *Event) String() string {
	return fmt.Sprintf("Event{%s, %s}", e.Hash().ShortString(), e.Parents.ShortString())
}