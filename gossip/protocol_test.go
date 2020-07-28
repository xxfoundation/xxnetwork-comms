package gossip

import (
	"crypto/rand"
	"gitlab.com/elixxir/primitives/id"
	"math"
	"sync"
	"testing"
)

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

// Auxiliar function to run all the tests of getPeers() result
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
	p := &Protocol {
		comms:            nil,
		fingerprints:     nil,
		fingerprintsLock: sync.RWMutex{},
		peers:            createListOfPeers(size, t),
		peersLock:        sync.RWMutex{},
		flags:            DefaultFlags(),
		receiver:         nil,
		verify:           nil,
		IsDefunct:        false,
	}

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