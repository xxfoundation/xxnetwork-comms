package gossip

import (
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"testing"
	"time"
)

// Basic test on manager creation
func TestNewManager(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	if m.buffer == nil || m.protocols == nil {
		t.Error("Failed to initialize all fields properly")
	}
}

// Happy path test for adding new gossip protocol
func TestManager_NewGossip(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())

	r := func(msg *GossipMsg) error {
		return nil
	}
	v := func(msg *GossipMsg) error {
		return nil
	}
	m.NewGossip("test", DefaultProtocolFlags(), r, v, []*id.ID{})

	if len(m.protocols) != 1 {
		t.Errorf("Failed to add protocol")
	}
}

// Happy path test for new gossip protocol with messages in buffer for that tag
func TestManager_NewGossip_WithBuffer(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	m.buffer["test"] = &MessageRecord{
		Timestamp: time.Time{},
		Messages:  []*GossipMsg{{Tag: "testmsg"}},
	}

	originalBufferLen := len(m.buffer)

	var received bool
	r := func(msg *GossipMsg) error {
		received = true
		return nil
	}
	v := func(msg *GossipMsg) error {
		return nil
	}
	m.NewGossip("test", DefaultProtocolFlags(), r, v, []*id.ID{})

	if len(m.protocols) != 1 {
		t.Errorf("Failed to add protocol")
	}

	if !received {
		t.Errorf("Did not receive message in buffer")
	}

	if originalBufferLen-len(m.buffer) != 1 {
		t.Errorf("Did not clear buffer after reception")
	}
}

// Basic unit test for getting a protocol
func TestManager_Get(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	m.protocols = map[string]*Protocol{"test": {}}

	_, ok := m.Get("test")
	if !ok {
		t.Error("Did not get protocol")
	}
}

// Basic unit test for deleting a protocol
func TestManager_Delete(t *testing.T) {
	m := NewManager(&connect.ProtoComms{}, DefaultManagerFlags())
	m.protocols = map[string]*Protocol{"test": {}}

	m.Delete("test")
	if len(m.protocols) != 0 {
		t.Error("Failed to delete protocol")
	}
}

// Test buffer monitor thread operation
func TestManager_BufferMonitor(t *testing.T) {
	flags := DefaultManagerFlags()
	flags.BufferExpirationTime = 3 * time.Second
	flags.MonitorThreadFrequency = 3 * time.Second
	m := NewManager(&connect.ProtoComms{}, flags)
	m.buffer["test"] = &MessageRecord{
		Timestamp: time.Now(),
		Messages:  nil,
	}
	origLen := len(m.buffer)
	time.Sleep(4 * time.Second)
	if origLen-len(m.buffer) != 1 {
		t.Errorf("Failed to clear buffer after duration expired")
	}
}
