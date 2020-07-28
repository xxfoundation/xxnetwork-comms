package gossip

import (
	"context"
	"crypto/md5"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
	"sync"
)

// Defines the type of Gossip message fingerprints
// hash(tag, origin, payload, signature)
type Fingerprint [16]byte

// Gossip-related configuration flag
type Flags struct {
	FanOut                  uint8  // Default = 0
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
	fingerprints     map[Fingerprint]struct{}
	fingerprintsLock sync.RWMutex

	// Thread-safe list of peers for the Protocol
	peers     []*id.ID
	peersLock sync.RWMutex

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
func (p *Protocol) receive(msg *GossipMsg) error {
	p.fingerprintsLock.Lock()
	defer p.fingerprintsLock.Unlock()

	// Get fingerprint
	preSum := append([]byte(msg.Tag), msg.Origin...)
	preSum = append(preSum, msg.Payload...)
	preSum = append(preSum, msg.Signature...)
	fingerprint := NewFingerprint(md5.Sum(preSum))
	if _, ok := p.fingerprints[fingerprint]; ok {
		return nil
	}

	err := p.verify(msg, msg.Signature)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify gossip message")
	}
	err = p.receiver(msg)
	if err != nil {
		return errors.WithMessage(err, "Failed to receive gossip message")
	}

	p.fingerprints[fingerprint] = struct{}{}
	return nil
}

// Adds a peer by ID to the Gossip protocol
func (p *Protocol) AddGossipPeer(id *id.ID) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	// Confirm we have a host matching this ID
	_, ok := p.comms.GetHost(id)
	if !ok {
		return errors.Errorf("Could not retreive host for ID %s", id)
	}

	p.peers = append(p.peers, id)
	return nil
}

// Builds and sends a GossipMsg
func (p *Protocol) Gossip(msg *GossipMsg) error {
	p.peersLock.RLock()
	defer p.peersLock.RUnlock()

	// Internal helper to send the input gossip msg to a given id
	sendFunc := func(id *id.ID) error {
		h, ok := p.comms.GetHost(id)
		if !ok {
			return errors.Errorf("Failed to get host with ID %s", id)
		}
		f := func(conn *grpc.ClientConn) (*any.Any, error) {
			gossipClient := NewGossipClient(conn)
			ack, err := gossipClient.Endpoint(context.Background(), msg)
			if err != nil {
				return nil, errors.WithMessage(err, "Failed to send message")
			}
			return ptypes.MarshalAny(ack)
		}
		_, err := p.comms.Send(h, f)
		if err != nil {
			return errors.WithMessagef(err, "Failed to send to host %s", h.String())
		}
		return nil
	}

	// Get list of peers to send message to
	peers, err := p.getPeers()
	if err != nil {
		return errors.WithMessage(err, "Failed to get peers for sending")
	}

	// Send message to each peer
	errCount := 0
	errs := errors.New("Failed to send message to some peers...")
	for _, p := range peers {
		sendErr := sendFunc(p) // TODO: Should this happen in a gofunc?
		if sendErr != nil {
			errs = errors.WithMessagef(err, "Failed to send to ID %s", p)
			errCount++
		}
	}
	if errCount > 0 {
		return errs
	} else {
		return nil
	}
}

// Performs returns which peers to send the GossipMsg to
func (p *Protocol) getPeers() ([]*id.ID, error) {
	return nil, nil
}

// NewFingerprint creates a new fingerprint from a byte slice
func NewFingerprint(data [16]byte) Fingerprint {
	fp := Fingerprint{}
	copy(fp[:], data[:])
	return fp
}
