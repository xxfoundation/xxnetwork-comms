package gossip

import (
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m.buffer == nil || m.protocolLock == nil || m.protocols == nil || m.bufferLock == nil {
		t.Error("Failed to initialize all fields properly")
	}
}

func TestManager_NewGossip(t *testing.T) {
	m := NewManager()

	r := func(msg *GossipMsg) error {
		return nil
	}
	v := func(msg *GossipMsg, b []byte) error {
		return nil
	}
	m.NewGossip(nil, "test", DefaultFlags(), r, v, []*id.ID{})

	if len(m.protocols) != 1 {
		t.Errorf("Failed to add protocol")
	}
}

func TestManager_Get(t *testing.T) {
	m := NewManager()
	m.protocols = map[string]*Protocol{"test": {}}

	_, ok := m.Get("test")
	if !ok {
		t.Error("Did not get protocol")
	}
}

func TestManager_Delete(t *testing.T) {
	m := NewManager()
	m.protocols = map[string]*Protocol{"test": {}}

	m.Delete("test")
	if len(m.protocols) != 0 {
		t.Error("Failed to delete protocol")
	}
}

func TestManager_Defunct(t *testing.T) {
	m := NewManager()
	m.protocols = map[string]*Protocol{"test": {IsDefunct: false}}

	m.Defunct("test")
	if !m.protocols["test"].IsDefunct {
		t.Error("Failed to mark protocol as defunct")
	}
}

func TestManager_GRPCRegister(t *testing.T) {
	// How do i test this?
}
