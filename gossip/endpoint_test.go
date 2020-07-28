package gossip

import (
	"context"
	"gitlab.com/elixxir/primitives/id"
	"testing"
	"time"
)

func TestComms_Endpoint_toProtocol(t *testing.T) {
	gossipComms := &Comms{
		ProtoComms: nil,
		Manager:    NewManager(),
	}

	var received bool
	r := func(msg *GossipMsg) error {
		received = true
		return nil
	}
	v := func(*GossipMsg, []byte) error {
		return nil
	}
	gossipComms.NewGossip(nil, "test", DefaultProtocolFlags(), r, v,
		[]*id.ID{id.NewIdFromString("test", id.Node, t)})

	_, err := gossipComms.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send: %+v", err)
	}

	if !received {
		t.Errorf("Didn't receive message in protocol")
	}
}

func TestComms_Endpoint_toNewBuffer(t *testing.T) {
	gossipComms := &Comms{
		ProtoComms: nil,
		Manager:    NewManager(),
	}
	_, err := gossipComms.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := gossipComms.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 1 {
		t.Errorf("Did not add message to buffer")
	}
}

func TestComms_Endpoint_toExistingBuffer(t *testing.T) {
	gossipComms := &Comms{
		ProtoComms: nil,
		Manager:    NewManager(),
	}
	now := time.Now()
	gossipComms.buffer["test"] = &MessageRecord{
		Timestamp: &now,
		Messages:  []*GossipMsg{{Tag: "test"}},
	}
	_, err := gossipComms.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := gossipComms.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 2 {
		t.Errorf("Did not add message to buffer")
	}
}

func TestComms_Stream(t *testing.T) {
	// TODO: Implement test once streaming is enabled
}
