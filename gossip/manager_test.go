package gossip

import (
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager(DefaultManagerFlags())
	if m.buffer == nil || m.protocols == nil {
		t.Error("Failed to initialize all fields properly")
	}
}

func TestManager_NewGossip(t *testing.T) {
	m := NewManager(DefaultManagerFlags())

	r := func(msg *GossipMsg) error {
		return nil
	}
	v := func(msg *GossipMsg, b []byte) error {
		return nil
	}
	m.NewGossip(nil, "test", DefaultProtocolFlags(), r, v, []*id.ID{})

	if len(m.protocols) != 1 {
		t.Errorf("Failed to add protocol")
	}
}

func TestManager_Get(t *testing.T) {
	m := NewManager(DefaultManagerFlags())
	m.protocols = map[string]*Protocol{"test": {}}

	_, ok := m.Get("test")
	if !ok {
		t.Error("Did not get protocol")
	}
}

func TestManager_Delete(t *testing.T) {
	m := NewManager(DefaultManagerFlags())
	m.protocols = map[string]*Protocol{"test": {}}

	m.Delete("test")
	if len(m.protocols) != 0 {
		t.Error("Failed to delete protocol")
	}
}

func TestManager_GRPCRegister(t *testing.T) {
	// How do i test this?
}
