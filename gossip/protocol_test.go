package gossip

import (
	"crypto/rand"
	"errors"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"math"
	"sync"
	"testing"
)

//====================================================================================================================//
// Basic unit tests of protocol methods
//====================================================================================================================//

// Test functionality of AddGossipPeer
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

// Happy path test for sending a gossip message
func TestProtocol_Gossip(t *testing.T) {
	// TODO: how should this be tested when we don't have getpeers implementation
	p := setup(t)
	p.flags.FanOut = 2

	testHostID := id.NewIdFromString("testhost", id.Node, t)
	testHostID2 := id.NewIdFromString("testhost2", id.Node, t)
	p.peers = []*id.ID{testHostID, testHostID2}

	_, errs := p.Gossip(&GossipMsg{
		Tag:       "test",
		Origin:    nil,
		Payload:   nil,
		Signature: nil,
	})
	// Since the hosts are fake, we will receive two errors (one for each failed attempt)
	if len(errs) != 2 {
		t.Errorf("Failed to send gossip messae: %+v", errs)
	}
}

// Happy path test for receive method
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

// Basic unit test for Defunct function on a protocol
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

// Setup a gossip protocol for basic testing - fields can be overridden as needed in tests
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
		fingerprints:     map[Fingerprint]*uint64{},
		fingerprintsLock: sync.RWMutex{},
		peers:            []*id.ID{},
		peersLock:        sync.RWMutex{},
		flags:            ProtocolFlags{},
		receiver:         r,
		verify:           v,
		IsDefunct:        false,
	}
}

// Test uniqueness of fingerprint function
func TestGetFingerprint(t *testing.T) {
	message1 := &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}
	message2 := &GossipMsg{
		Tag:       "watermelon",
		Origin:    []byte("apple"),
		Payload:   []byte("quantum physics"),
		Signature: []byte("kurask"),
	}
	message3 := &GossipMsg{
		Tag:       "tesv",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}
	message4 := &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	}
	f1 := GetFingerprint(message1)
	f2 := GetFingerprint(message2)
	f3 := GetFingerprint(message3)
	f4 := GetFingerprint(message4)
	if f1 == f2 || f2 == f3 || f1 == f3 {
		t.Errorf("Fingerprints formed from unique messages are not unique")
	}
	if f1 != f4 {
		t.Errorf("Fingerprints formed from identical messages are not identical")
	}
}

//====================================================================================================================//
// Testing Reader of getPeers()
//====================================================================================================================//
type ErrReader struct{}

func (r *ErrReader) Read(p []byte) (n int, err error) {
	return len(p), errors.New("TEST")
}

func testGetPeersReader(p *Protocol, t *testing.T) {
	// Test getPeers() with real reader
	p.crand = rand.Reader
	_, err := p.getPeers()
	if err != nil {
		t.Errorf("[Test Real Reader] getPeers() error = %v", err)
	}

	// Test getPeers() with error reader
	p.crand = &ErrReader{}
	_, err = p.getPeers()
	if err == nil {
		t.Errorf("[Test Error Reader] getPeers() error = %v", err)
	}
}

//====================================================================================================================//
// Testing Result of getPeers()
//====================================================================================================================//

// Generates a byte slice of the specified length containing random numbers.
func newRandomBytes(length int, t *testing.T) []byte {
	// Create new byte slice of the correct size
	idBytes := make([]byte, length)

	// Create random bytes
	_, err := rand.Read(idBytes)
	if err != nil {
		t.Fatalf("Failed to generate random bytes: %v", err)
	}

	return idBytes
}

// Create a random list of peers
func createListOfPeers(size int, t *testing.T) []*id.ID {
	list := make([]*id.ID, size)
	for i := 0; i < size; i++ {
		randomBytes := newRandomBytes(id.ArrIDLen, t)
		list[i] = id.NewIdFromBytes(randomBytes, t)
	}
	return list
}

// Check if sub-list of peers is contained in peers full list
func containedIn(result, peers []*id.ID) bool {
	for i := 0; i < len(result); i++ {
		if !contains(peers, result[i]) {
			return false
		}
	}
	return true
}

// Check if peer is contained in a peers list
func contains(list []*id.ID, peer *id.ID) bool {
	for _, a := range list {
		if a == peer {
			return true
		}
	}
	return false
}

// Check if resulted peers list has duplicates
func hasDuplicates(list []*id.ID) bool {
	auxMap := make(map[*id.ID]struct{})
	for _, a := range list {
		if _, ok := auxMap[a]; ok {
			return true
		}
		auxMap[a] = struct{}{}
	}
	return false
}

// Test result returned from getPeers()
func testGetPeers(p *Protocol, t *testing.T) {
	// Run tested method
	list, err := p.getPeers()
	if err != nil {
		t.Errorf("getPeers() error = %v", err)
		return
	}

	// Debug results
	t.Logf("\tResult (%d) = %v\n", len(list), list)

	// Set the expected values
	expectedSize := int(p.flags.FanOut)
	if expectedSize == 0 {
		expectedSize = int(math.Ceil(math.Sqrt(float64(len(p.peers)))))
	} else if expectedSize > len(p.peers) {
		expectedSize = len(p.peers)
	}

	// Test resulted list as duplicated peers
	if hasDuplicates(list) {
		t.Errorf("getPeers() returned list with duplicated peers! Fanout = %d", p.flags.FanOut)
	}

	// Test if resulted list is contained in peers list
	if !containedIn(list, p.peers) {
		t.Errorf("getPeers() returned list with peers not part of p.peers! Fanout = %d", p.flags.FanOut)
	}

	// Test if resulted list size matches fanout/expectedSize
	if len(list) != expectedSize {
		t.Errorf("getPeers() returned unexpected list size! Fanout = %d", p.flags.FanOut)
	}
}

func TestProtocol_getPeers(t *testing.T) {
	// Initialize Default Variables
	size := 10
	p := &Protocol{
		comms:            nil,
		fingerprints:     nil,
		fingerprintsLock: sync.RWMutex{},
		peers:            createListOfPeers(size, t),
		peersLock:        sync.RWMutex{},
		flags:            DefaultProtocolFlags(),
		receiver:         nil,
		verify:           nil,
		IsDefunct:        false,
		crand:            rand.Reader,
	}

	// Test Reader
	testGetPeersReader(p, t)

	// Restore the Reader
	p.crand = rand.Reader

	t.Logf("Initial Variables:")
	t.Logf("\tPeers (%d) = %v", len(p.peers), p.peers)

	// Case 1: fanout = 0, should return a set of peers contained in peers []*id.ID and with size = sqrt(len(peers))
	p.flags.FanOut = 0
	t.Logf("Case 1: fanout = %d", p.flags.FanOut)
	testGetPeers(p, t)

	// Case 2: fanout != 0, should return a set of peers contained in peers []*id.ID and with size = fanout
	p.flags.FanOut = uint8(size / 2.0)
	t.Logf("Case 2: fanout (%d) > 0", p.flags.FanOut)
	testGetPeers(p, t)

	// Case 3: fanout > size, should return a set of peers contained in peers []*id.ID and with size = fanout
	p.flags.FanOut = uint8(size + 2)
	t.Logf("Case 3: fanout (%d) > size (%d)", p.flags.FanOut, size)
	testGetPeers(p, t)
}
