///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Gossip protocol definition

package gossip

import (
	"context"
	"crypto/md5"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/shuffle"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"io"
	"math"
	"sync"
	"sync/atomic"
)

// Defines the type of Gossip message fingerprints
// hash(tag, origin, payload, signature)
type Fingerprint [16]byte

const minimumPeers = 20

// NewFingerprint creates a new fingerprint from a byte slice
func NewFingerprint(data [16]byte) Fingerprint {
	fp := Fingerprint{}
	copy(fp[:], data[:])
	return fp
}

// Obtain the fingerprint of the GossipMsg
func GetFingerprint(msg *GossipMsg) Fingerprint {
	preSum := append([]byte(msg.Tag), msg.Origin...)
	preSum = append(preSum, msg.Payload...)
	preSum = append(preSum, msg.Signature...)
	fingerprint := NewFingerprint(md5.Sum(preSum))
	return fingerprint
}

// Returns the data of a GossipMsg excluding the Signature as bytes
func Marshal(msg *GossipMsg) []byte {
	data := append([]byte(msg.Tag), msg.Origin...)
	return append(data, msg.Payload...)
}

// Gossip-related configuration flag
type ProtocolFlags struct {
	FanOut                  uint8  // Default = 0
	MaxRecordedFingerprints uint64 // Default = 10000000
	MaximumReSends          uint64 // Default = 3
	NumParallelSends        uint8  // Default = 5
}

// Returns a ProtocolFlags object with all flags set to their defaults
func DefaultProtocolFlags() ProtocolFlags {
	return ProtocolFlags{
		FanOut:                  0,
		MaxRecordedFingerprints: 10000000,
		MaximumReSends:          3,
		NumParallelSends:        30,
	}
}

// Generic interface representing various Gossip protocols
type Protocol struct {
	comms *connect.ProtoComms

	// Thread-safe record of all Gossip messages for this Protocol
	// NOTE: Must avoid unlimited growth
	fingerprints     map[Fingerprint]*uint64
	fingerprintsLock sync.RWMutex

	// Thread-safe list of peers for the Protocol
	peers     []*id.ID
	peersLock sync.RWMutex

	// Stores the Gossip-related configuration flags
	flags ProtocolFlags

	// Handler function for GossipMsg Reception
	receiver Receiver

	// Verifier function for GossipMsg signatures
	verify SignatureVerification

	// Marks a Protocol as Defunct such that it will ignore new messages
	IsDefunct   bool
	defunctLock sync.Mutex

	// Random Reader
	crand io.Reader

	// worker pool channel for sending
	sendWorkers chan sendInstructions
}

type sendInstructions struct {
	sendFunc   func(id *id.ID) error
	peer       *id.ID
	errChannel chan error
	wait       *sync.WaitGroup
}

// Marks a Protocol as Defunct such that it will ignore new messages
func (p *Protocol) Defunct() {
	p.defunctLock.Lock()
	p.IsDefunct = true
	p.defunctLock.Unlock()
}

// Receive a Gossip Message and check fingerprints map
// (if unique calls GossipSignatureVerify -> Receiver)
func (p *Protocol) receive(msg *GossipMsg) error {
	var err error

	// Check fingerprint of the message against our record
	fingerprint := GetFingerprint(msg)
	p.fingerprintsLock.RLock()
	numSendsPrt, ok := p.fingerprints[fingerprint]
	p.fingerprintsLock.RUnlock()
	//if there is no record of receiving the fingerprint, process it as new
	if !ok {
		err = p.verify(msg, nil)
		if err != nil {
			return errors.WithMessage(err, "Failed to verify gossip message")
		}
		p.fingerprintsLock.Lock()
		// redo the check if the fingerprint exists to ensure it was not added
		// between the previous check and taking the lock
		numSendsPrt, ok = p.fingerprints[fingerprint]
		if !ok {
			numSendsLocal := uint64(0)
			p.fingerprints[fingerprint] = &numSendsLocal
			p.fingerprintsLock.Unlock()
			err = p.receiver(msg)
			if err != nil {
				return errors.WithMessage(err, "Failed to receive gossip message")
			}
		} else {
			p.fingerprintsLock.Unlock()
		}
	}

	// Increment the number of sends for this fingerprint
	numSends := uint64(0)
	if ok {
		numSends = atomic.AddUint64(numSendsPrt, 1)
	}

	if numSends < p.flags.MaximumReSends {
		// Since gossip propagates the message across a potentially large message, we don't want this to block
		go func() {
			numPeers, errs := p.Gossip(msg)
			if len(errs) != 0 {
				jww.TRACE.Print(errors.Errorf("Failed to gossip message to %d of %d peers", len(errs), numPeers))
			}
		}()
	}

	return nil
}

// Adds a peer by ID to the Gossip protocol
func (p *Protocol) AddGossipPeer(id *id.ID) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	// Confirm we have a host matching this ID
	// Because hosts can be removed, this CANNOT ensure the host still exists when an ID is used
	_, ok := p.comms.GetHost(id)
	if !ok {
		return errors.Errorf("Could not retrieve host for ID %s", id)
	}

	p.peers = append(p.peers, id)
	return nil
}

// Returns all peers in the gossip. Primarily for debugging.
func (p *Protocol) GetPeers() []*id.ID {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	peersCopy := make([]*id.ID, len(p.peers))

	for i, peer := range p.peers {
		peersCopy[i] = peer.DeepCopy()
	}

	return peersCopy
}

// Remove a peer by ID to the Gossip protocol
func (p *Protocol) RemoveGossipPeer(id *id.ID) error {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	for i, peer := range p.peers {
		if peer.Cmp(id) {
			p.peers = append(p.peers[:i], p.peers[i+1:]...)
			return nil
		}
	}

	return errors.Errorf("Could not remove peer for ID %s", id)
}

// Builds and sends a GossipMsg
func (p *Protocol) Gossip(msg *GossipMsg) (int, []error) {
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
		return 0, []error{errors.WithMessage(err, "Failed to get peers for sending")}
	}

	// Send message to each peer
	errCh := make(chan error, len(peers))
	wg := sync.WaitGroup{}
	wg.Add(len(peers))
	//send signals to the worker threads to do the sends
	for _, peer := range peers {
		p.sendWorkers <- sendInstructions{
			sendFunc:   sendFunc,
			peer:       peer.DeepCopy(),
			errChannel: errCh,
			wait:       &wg,
		}
	}

	//wait for sends to complete
	wg.Wait()

	//get any errors
	done := false
	var errs []error
	for !done {
		//pull from the channel until errors are exhausted
		select {
		case newErr := <-errCh:
			errs = append(errs, newErr)
		default:
			done = true
		}
	}

	if len(errs) > 0 {
		return len(peers), errs
	} else {
		return len(peers), nil
	}
}

// Performs returns which peers to send the GossipMsg to
func (p *Protocol) getPeers() ([]*id.ID, error) {
	p.peersLock.RLock()
	defer p.peersLock.RUnlock()

	// Check fanout
	size := len(p.peers)
	fanout := int(p.flags.FanOut)

	if p.flags.FanOut < 1 {
		fanout = int(math.Ceil(math.Sqrt(float64(size))))
	}
	if size <= fanout || size < minimumPeers {
		return p.peers, nil
	}

	// Compute seed
	seed := make([]byte, 32)
	_, err := p.crand.Read(seed)
	if err != nil {
		return nil, err
	}

	// Use Fisher Yates Shuffle
	out := make([]*id.ID, fanout)
	shuffled := shuffle.SeededShuffle(size, seed)
	for i := 0; i < fanout; i++ {
		out[i] = p.peers[shuffled[i]]
	}

	return out, nil
}
