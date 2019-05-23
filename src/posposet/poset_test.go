package posposet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/sorting"
)

func TestPoset(t *testing.T) {
	nodes, nodesEvents := GenEventsByNode(5, 99, 3)

	posets := make([]*Poset, len(nodes))
	inputs := make([]*EventStore, len(nodes))
	for i := 0; i < len(nodes); i++ {
		posets[i], _, inputs[i] = FakePoset(nodes)
	}

	t.Run("Multiple start", func(t *testing.T) {
		posets[0].Stop()
		posets[0].Start()
		posets[0].Start()
	})

	t.Run("Push unordered events", func(t *testing.T) {
		// first all events from one node
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {
				inputs[n].SetEvent(&e.Event)
				posets[n].PushEventSync(e.Hash())
			}
		}
		// second all events from others
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {
				for i := 0; i < len(posets); i++ {
					if i != n {
						inputs[i].SetEvent(&e.Event)
						posets[i].PushEventSync(e.Hash())
					}
				}
			}
		}
	})

	t.Run("All events in Store", func(t *testing.T) {
		for _, events := range nodesEvents {
			for _, e0 := range events {
				frame := posets[0].store.GetEventFrame(e0.Hash())
				if !assert.NotNil(t, frame, "Event is not in poset store") {
					return
				}
			}
		}
	})

	t.Run("Check consensus", func(t *testing.T) {
		asserts := assert.New(t)
		for i := 0; i < len(posets)-1; i++ {
			for j := i + 1; j < len(posets); j++ {
				p0, p1 := posets[i], posets[j]
				// compare blockchain
				if !asserts.Equal(p0.state.LastBlockN, p1.state.LastBlockN, "blocks count") {
					return
				}
				for b := uint64(1); b <= p0.state.LastBlockN; b++ {
					if !asserts.Equal(p0.store.GetBlock(b), p1.store.GetBlock(b), "block") {
						return
					}
				}
			}
		}
	})

	t.Run("Multiple stop", func(t *testing.T) {
		posets[0].Stop()
		posets[0].Stop()
	})
}

func TestPosetWithOrdering(t *testing.T) {
	nodes, nodesEvents := GenEventsByNode(5, 99, 3)

	posets := make([]*Poset, len(nodes))
	inputs := make([]*EventStore, len(nodes))
	order := make([]*sorting.Ordering, len(nodes))
	for i := 0; i < len(nodes); i++ {
		posets[i], _, inputs[i], order[i] = FakePosetWithOrdering(nodes)
	}

	t.Run("Multiple start", func(t *testing.T) {
		posets[0].Stop()
		posets[0].Start()
		posets[0].Start()
	})

	t.Run("Push unordered events", func(t *testing.T) {
		// first all events from one node
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {
				inputs[n].SetEvent(&e.Event)
				posets[n].PushEventWithOrdering(&e.Event, order[n])
			}
		}
		// second all events from others
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {

				for i := 0; i < len(posets); i++ {
					if i != n {
						inputs[i].SetEvent(&e.Event)
						posets[i].PushEventWithOrdering(&e.Event, order[i])
					}
				}
			}
		}
	})

	t.Run("All events in Store", func(t *testing.T) {
		for _, events := range nodesEvents {
			for _, e0 := range events {
				frame := posets[0].store.GetEventFrame(e0.Hash())
				if !assert.NotNil(t, frame, "Event is not in poset store") {
					return
				}
			}
		}
	})

	t.Run("Check consensus", func(t *testing.T) {
		asserts := assert.New(t)
		for i := 0; i < len(posets)-1; i++ {
			for j := i + 1; j < len(posets); j++ {
				p0, p1 := posets[i], posets[j]
				// compare blockchain
				if !asserts.Equal(p0.state.LastBlockN, p1.state.LastBlockN, "blocks count") {
					return
				}
				for b := uint64(1); b <= p0.state.LastBlockN; b++ {
					if !asserts.Equal(p0.store.GetBlock(b), p1.store.GetBlock(b), "block") {
						return
					}
				}
			}
		}
	})

	t.Run("Multiple stop", func(t *testing.T) {
		posets[0].Stop()
		posets[0].Stop()
	})
}

func TestOrdering(t *testing.T) {
	nodesCount, eventsCount, parentCount := 5, 99, 3
	nodes, nodesEvents := GenEventsByNode(nodesCount, eventsCount, parentCount)

	inputs := make([]sorting.Source, len(nodes))
	orders := make([]*sorting.Ordering, len(nodes))
	for i := 0; i < len(nodes); i++ {
		input := NewEventStore(nil, false)
		order := sorting.NewOrder(input, log, true)
		inputs[i], orders[i] = input, order
	}

	t.Run("Push unordered events", func(t *testing.T) {
		// first all events from one node
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {
				inputs[n].SetEvent(&e.Event)
				ch := orders[n].NewEventConsumer(&e.Event)
				orders[n].NewEventOrder(ch)
			}
		}

		// second all events from others
		for n := 0; n < len(nodes); n++ {
			events := nodesEvents[nodes[n]]
			for _, e := range events {

				for i := 0; i < len(inputs); i++ {
					if i != n {
						inputs[i].SetEvent(&e.Event)
						ch := orders[i].NewEventConsumer(&e.Event)
						orders[i].NewEventOrder(ch)
					}
				}
			}
		}
	})
}

/*
 * Poset's test methods:
 */

// PushEventSync takes event into processing. It's a sync version of Poset.PushEvent().
// Event order doesn't matter.
func (p *Poset) PushEventSync(e hash.Event) {
	event := p.GetEvent(e)
	p.onNewEvent(event)
}

func (p *Poset) PushEventWithOrdering(e *inter.Event, o *sorting.Ordering) {
	ch := o.NewEventConsumer(e, p.consensus)
	o.NewEventOrder(ch, p.store.GetEventFrame)
}
