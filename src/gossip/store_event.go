package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/src/inter/wire"
)

// TODO store separately
// GetEventHeader returns stored event header.
func (s *Store) GetEventHeader(h hash.Event) *inter.EventHeaderData {
	e := s.GetEvent(h)
	if e == nil {
		return nil
	}
	return &e.EventHeaderData
}

// SetEvent stores event.
func (s *Store) SetEvent(e *inter.Event) {
	key := e.Hash().Bytes()

	w := e.ToWire()

	s.set(s.table.Events, key, w)
}

// GetEvent returns stored event.
func (s *Store) GetEvent(h hash.Event) *inter.Event {
	w := s.GetWireEvent(h)
	return inter.WireToEvent(w)
}

// HasEvent returns true if event exists.
func (s *Store) HasEvent(h hash.Event) bool {
	return s.has(s.table.Events, h.Bytes())
}

// GetWireEvent returns stored event.
// Result is a ready gRPC message.
func (s *Store) GetWireEvent(h hash.Event) *wire.Event {
	key := h.Bytes()
	w, _ := s.get(s.table.Events, key, &wire.Event{}).(*wire.Event)
	if w == nil {
		return w
	}

	return w
}

// SetEventHash stores hash.
func (s *Store) SetEventHash(creator hash.Peer, sf idx.SuperFrame, seq idx.Event, hash hash.Event) {

	key := append(creator.Bytes(), sf.Bytes()...)
	key = append(key, seq.Bytes()...)

	if err := s.table.Hashes.Put(key, hash.Bytes()); err != nil {
		s.Fatal(err)
	}
}

// GetEventHash returns stored event hash.
func (s *Store) GetEventHash(creator hash.Peer, sf idx.SuperFrame, seq idx.Event) *hash.Event {
	key := append(creator.Bytes(), sf.Bytes()...)
	key = append(key, seq.Bytes()...)

	buf, err := s.table.Hashes.Get(key)
	if err != nil {
		s.Fatal(err)
	}
	if buf == nil {
		return nil
	}

	e := hash.BytesToEvent(buf)
	return &e
}