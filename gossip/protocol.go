package gossip

import (
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// Defines the type of Gossip message fingerprints
// hash(tag, origin, payload, signature)
type Fingerprint [32]byte

// Gossip-related configuration flag
type Flags struct {
	FanOut                  uint8
	MaxRecordedFingerprints uint64 // Default = 10000000
}

// Returns a Flags object with all flags set to their defaults
func DefaultFlags() Flags {
	return Flags{}
}

// Generic interface representing various Gossip protocols
type Protocol struct {
	comms *connect.ProtoComms

	// Thread-safe record of all Gossip messages for this Protocol
	// NOTE: Must avoid unlimited growth
	fingerprints map[Fingerprint]struct{}

	// Thread-safe list of peers for the Protocol
	peers []*id.ID

	// Stores the Gossip-related configuration flags
	flags Flags

	// Handler function for GossipMsg Reception
	receiver Receiver

	// Verifier function for GossipMsg signatures
	verify SignatureVerification

	// Marks a Protocol as Defunct such that it will ignore new messages
	IsDefunct bool
}

// Receive a Gossip Message and check fingerprints map
// (if unique calls GossipSignatureVerify -> Receiver)
func (p *Protocol) receive(msg *messages.GossipMsg) error {
	return nil
}

// Adds a peer by ID to the Gossip protocol
func (p *Protocol) AddGossipPeer(id *id.ID) error {
	return nil
}

// Builds and sends a GossipMsg
func (p *Protocol) Gossip(msg *messages.GossipMsg) error {
	return nil
}

// Performs returns which peers to send the GossipMsg to
func (p *Protocol) getPeers() ([]*id.ID, error) {
	return nil, nil
}
