package gossip

import (
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"sync"
	"testing"
)

func TestProtocol_AddGossipPeer(t *testing.T) {
	p := setup(t)
	testHostID := id.NewIdFromString("testhost", id.Node, t)
	_, err := p.comms.AddHost(testHostID, "0.0.0.0:420", nil, true, false)
	if err != nil {
		t.Errorf("Failed to add test host: %+v", err)
	}
	err = p.AddGossipPeer(testHostID)
	if err != nil {
		t.Errorf("Failed to add gossip peer: %+v", err)
	}
	if len(p.peers) == 0 {
		t.Errorf("Did not properly add gossip peer")
	}
}

func TestProtocol_Gossip(t *testing.T) {
	// TODO: how should this be tested when we don't have getpeers implementation
	p := setup(t)
	_, errs := p.Gossip(&GossipMsg{
		Tag:       "test",
		Origin:    nil,
		Payload:   nil,
		Signature: nil,
	})
	if len(errs) != 0 {
		t.Errorf("Failed to send gossip message: %+v", errs)
	}
}

func TestProtocol_receive(t *testing.T) {
	p := setup(t)
	r := func(msg *GossipMsg) error {
		return nil
	}
	p.receiver = r

	message1 := &GossipMsg{
		Tag:       "Message1",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}

	message2 := &GossipMsg{
		Tag:       "Message2",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}

	err := p.receive(message1)
	if err != nil {
		t.Errorf("Failed to receive message1: %+v", err)
	}
	if len(p.fingerprints) != 1 {
		t.Errorf("Did not add message1 fingerprint to array")
	}

	err = p.receive(message2)
	if err != nil {
		t.Errorf("Failed to receive message2: %+v", err)
	}
	if len(p.fingerprints) != 2 {
		t.Errorf("Did not add message2 fingerprint to array")
	}

	err = p.receive(message1)
	if err != nil {
		t.Errorf("Failed to receive duplicate of message1: %+v", err)
	}
	if len(p.fingerprints) != 2 {
		t.Errorf("Fingerprint of duplicate message was added to array")
	}
}

func TestProtocol_Defunct(t *testing.T) {
	p := Protocol{
		comms:            nil,
		fingerprints:     nil,
		fingerprintsLock: sync.RWMutex{},
		peers:            nil,
		peersLock:        sync.RWMutex{},
		flags:            ProtocolFlags{},
		receiver:         nil,
		verify:           nil,
		IsDefunct:        false,
		defunctLock:      sync.Mutex{},
	}

	p.Defunct()
	if !p.IsDefunct {
		t.Error("Failed to mark protocol as defunct")
	}
}

func setup(t *testing.T) *Protocol {
	r := func(msg *GossipMsg) error {
		return nil
	}
	v := func(msg *GossipMsg, b []byte) error {
		return nil
	}
	c := &connect.ProtoComms{
		Manager: connect.Manager{},
	}
	return &Protocol{
		comms:            c,
		fingerprints:     map[Fingerprint]uint64{},
		fingerprintsLock: sync.RWMutex{},
		peers:            nil,
		peersLock:        sync.RWMutex{},
		flags:            ProtocolFlags{},
		receiver:         r,
		verify:           v,
		IsDefunct:        false,
	}
}
