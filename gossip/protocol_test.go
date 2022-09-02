////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gossip

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/primitives/id"
	"math"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ====================================================================================================================//
// Basic unit tests of protocol methods
// ====================================================================================================================//

// Test functionality of AddGossipPeer
func TestProtocol_AddGossipPeer(t *testing.T) {
	p := setup(t)
	testHostID := id.NewIdFromString("testhost", id.Node, t)
	_, err := p.comms.AddHost(testHostID, "0.0.0.0:420", nil, connect.GetDefaultHostParams())
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

// Happy path
func TestProtocol_RemoveGossipPeer(t *testing.T) {
	p := setup(t)
	testHostID := id.NewIdFromString("testhost", id.Node, t)
	p.peers = append(p.peers, testHostID)
	if len(p.peers) != 1 {
		t.Errorf("Expected to add gossip peer")
	}
	err := p.RemoveGossipPeer(testHostID)
	if err != nil {
		t.Errorf("Unable to remove gossip peer: %+v", err)
	}
	if len(p.peers) != 0 {
		t.Errorf("Expected to remove gossip peer")
	}
}

// Test functionality of GetPeers
func TestProtocol_GetPeers(t *testing.T) {
	p := setup(t)
	testHostID := id.NewIdFromString("testhost", id.Node, t)
	_, err := p.comms.AddHost(testHostID, "0.0.0.0:420", nil, connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Failed to add test host: %+v", err)
	}
	err = p.AddGossipPeer(testHostID)
	if err != nil {
		t.Errorf("Failed to add gossip peer: %+v", err)
	}
	if len(p.GetPeers()) == 0 {
		t.Errorf("Did not properly add gossip peer")
	}
}

// Error path
func TestProtocol_RemoveGossipPeerError(t *testing.T) {
	p := setup(t)
	testHostID := id.NewIdFromString("testhost", id.Node, t)
	testHostID2 := id.NewIdFromString("testhost2", id.Node, t)
	p.peers = append(p.peers, testHostID)
	if len(p.peers) != 1 {
		t.Errorf("Expected to add gossip peer")
	}
	err := p.RemoveGossipPeer(testHostID2)
	if err == nil {
		t.Errorf("Improperly removed gossip peer: %+v", err)
	}
	if len(p.peers) != 1 {
		t.Errorf("Expected NOT to remove gossip peer")
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
		t.Errorf("Did not add message1 fingerprint to array: %d", len(p.fingerprints))
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

// Happy path test for receive method
func TestProtocol_receive_oldMessage(t *testing.T) {
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
		Timestamp: time.Now().Add(-time.Minute * 11).UnixNano(),
	}

	err := p.receive(message1)
	if err != nil {
		t.Errorf("Failed to receive message1: %+v", err)
	}
	if len(p.fingerprints) != 1 {
		t.Errorf("Did not add message1 fingerprint to array")
	}
	fingerprint := p.fingerprinter(message1)
	if atomic.LoadUint64(p.fingerprints[fingerprint]) != 0 {
		t.Error("Should not have gossiped old message")
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
	v := func(msg *GossipMsg, smth []byte) error {
		return nil
	}
	c := &connect.ProtoComms{
		Manager: connect.NewManagerTesting(t),
	}

	flags := DefaultProtocolFlags()
	p := &Protocol{
		comms:            c,
		fingerprints:     map[Fingerprint]*uint64{},
		oldFingerprints:  map[Fingerprint]*uint64{},
		fingerprintsLock: sync.RWMutex{},
		peers:            []*id.ID{},
		peersLock:        sync.RWMutex{},
		flags:            flags,
		receiver:         r,
		verify:           v,
		IsDefunct:        false,
		sendWorkers:      make(chan sendInstructions, 100*flags.NumParallelSends),
		fingerprinter: func(msg *GossipMsg) Fingerprint {
			return getFingerprint(msg)
		},
	}

	launchSendWorkers(flags.NumParallelSends, p.sendWorkers)
	return p
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
	f1 := getFingerprint(message1)
	f2 := getFingerprint(message2)
	f3 := getFingerprint(message3)
	f4 := getFingerprint(message4)
	if f1 == f2 || f2 == f3 || f1 == f3 {
		t.Errorf("Fingerprints formed from unique messages are not unique")
	}
	if f1 != f4 {
		t.Errorf("Fingerprints formed from identical messages are not identical")
	}
}

// ====================================================================================================================//
// Testing Reader of getPeers()
// ====================================================================================================================//
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

// ====================================================================================================================//
// Testing Result of getPeers()
// ====================================================================================================================//

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
		t.Errorf("getPeers() returned unexpected list size! Fanout = %d. "+
			"\n\tExpected: %d \n\tReceived: %d", p.flags.FanOut, expectedSize, len(list))
	}
}

func TestProtocol_getPeers(t *testing.T) {
	// Initialize Default Variables
	size := 25
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

// Show that the gossip protocol reliably shares data with all nodes
// Show that data that's valid propagates, and that data that's not valid doesn't
func TestGossipNodes(t *testing.T) {
	// Have some nodes
	numNodes := 50
	portOffset := 24671
	managers := make([]*Manager, 0, numNodes)
	nodes := make([]*id.ID, 0, numNodes)
	ports := make([]string, 0, numNodes)
	// atomic counter for gossip reception
	numReceived := uint64(0)
	listeners := make([]net.Listener, 0, numNodes)

	validSig := []byte("valid signature")
	// unsure if this is causing my problems
	certPEM := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	keyPEM := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
	// start comm servers
	for i := 0; i < numNodes; i++ {
		port := strconv.FormatUint(uint64(portOffset+i), 10)
		node := id.NewIdFromUInt(uint64(i), id.Node, t)
		nodes = append(nodes, node)
		ports = append(ports, port)
		pc, listen, err := connect.StartCommServer(node, "0.0.0.0:"+port, certPEM, keyPEM, nil)
		listeners = append(listeners, listen)

		// Do I need to add other hosts before calling NewManager?
		if err != nil {
			// Not startin' a node? that's a paddlin'
			t.Fatalf("error starting node %v on port %v w/ err %v", i, port, err)
		}
		// each server has one gossip manager
		manager := NewManager(pc, DefaultManagerFlags())
		managers = append(managers, manager)

		go func() {
			// start serving
			RegisterGossipServer(pc.LocalServer, manager)
			err := pc.LocalServer.Serve(listen)
			if err != nil {
				t.Fatal(err)
			}
		}()
	}

	// Have each node add all other nodes as peers
	for i := 0; i < numNodes; i++ {

		peers := make([]*id.ID, 0, numNodes)
		for j := 0; j < numNodes; j++ {
			if i != j {
				peers = append(peers, nodes[j])
				params := connect.GetDefaultHostParams()
				params.AuthEnabled = false
				params.MaxRetries = 3
				_, err := managers[i].comms.AddHost(nodes[j], "127.0.0.1:"+ports[j], certPEM, params)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
		protoFlags := DefaultProtocolFlags()
		protoFlags.NumParallelSends = 5
		protoFlags.MaxGossipAge = 30 * time.Second
		// Initialize a test gossip protocol on all the servers
		managers[i].NewGossip("test", protoFlags, func(msg *GossipMsg) error {
			// receive func
			atomic.AddUint64(&numReceived, 1)
			return nil
		}, func(msg *GossipMsg, something []byte) error {
			// check sig func
			if !bytes.Equal(validSig, msg.Signature) {
				return errors.New("verification error")
			}
			return nil
		}, peers)
		protocol, ok := managers[i].Get("test")
		if ok {
			// crand should be populated when calling NewGossip, right?
			// but, it isn't
			protocol.crand = rand.Reader
		}
	}
	// in case we need time for servers to come up
	time.Sleep(1 * time.Second)

	// Send some gossips!
	numToSend := 64
	for numSent := 0; numSent < numToSend; numSent++ {
		nodeIndex := numSent % numNodes
		protocol, ok := managers[nodeIndex].Get("test")
		if !ok {
			t.Fatalf("manager %d should have had a test protocol", nodeIndex)
		}
		payload := make([]byte, 8)
		binary.BigEndian.PutUint64(payload, uint64(numSent))
		_, errs := protocol.Gossip(&GossipMsg{
			Tag:       "test",
			Origin:    nodes[nodeIndex].Bytes(),
			Payload:   payload,
			Signature: validSig,
		})
		for i := range errs {
			if errs[i] != nil {
				t.Error(errs)
			}
		}
	}

	// wait for the received messages to get counted
	time.Sleep(1 * time.Second)
	numReceivedAfterABit := atomic.LoadUint64(&numReceived)
	if numReceivedAfterABit != uint64((numNodes-1)*numToSend) {
		t.Errorf("Received %d messages. Expected %d messages", numReceivedAfterABit, (numNodes-1)*numToSend)
	}
}
