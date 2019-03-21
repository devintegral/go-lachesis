package posnode

import (
	"crypto/ecdsa"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

// Node is a Lachesis node implementation.
type Node struct {
	ID    common.Hash
	key   *ecdsa.PrivateKey
	pub   *ecdsa.PublicKey
	store *Store

	peerDialer Dialer

	consensus Consensus

	service
	gossip
	discovery
}

// New creates node.
func New(key *ecdsa.PrivateKey, s *Store, c Consensus, peerDialer Dialer) *Node {
	return &Node{
		ID:    common.BytesToHash(crypto.FromECDSAPub(&key.PublicKey)),
		key:   key,
		pub:   &key.PublicKey,
		store: s,

		peerDialer: peerDialer,

		consensus: c,
	}
}

// Shutdown stops node.
func (n *Node) Shutdown() {

}