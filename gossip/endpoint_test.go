////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gossip

import (
	"context"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelDebug)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

// Test endpoint when manager has a protocol
func TestManager_Endpoint_toProtocol(t *testing.T) {
	pc := &connect.ProtoComms{
		Manager: connect.NewManagerTesting(t),
	}

	m := NewManager(pc, DefaultManagerFlags())

	received := make(chan bool)
	r := func(msg *GossipMsg) error {
		received <- true
		return nil
	}
	v := func(*GossipMsg, []byte) error {
		return nil
	}
	m.NewGossip("test", DefaultProtocolFlags(), r, v,
		[]*id.ID{id.NewIdFromString("test", id.Node, t)})

	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send: %+v", err)
	}

	select {
	case <-received:

	case <-time.NewTimer(5 * time.Millisecond).C:
		t.Errorf("Didn't receive message in protocol")
	}
}

// Test endpoint function when there is no protocol and no buffer record
func TestManager_Endpoint_toNewBuffer(t *testing.T) {
	pc := &connect.ProtoComms{
		Manager: connect.NewManagerTesting(t),
	}

	m := NewManager(pc, DefaultManagerFlags())
	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := m.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 1 {
		t.Errorf("Did not add message to buffer")
	}
}

// Test endpoint function when there is no protocol, but an existing buffer
func TestComms_Endpoint_toExistingBuffer(t *testing.T) {
	pc := &connect.ProtoComms{
		Manager: connect.NewManagerTesting(t),
	}

	m := NewManager(pc, DefaultManagerFlags())
	now := time.Now()
	m.buffer["test"] = &MessageRecord{
		Timestamp: now,
		Messages:  []*GossipMsg{{Tag: "test"}},
	}
	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	r, ok := m.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(r.Messages) != 2 {
		t.Errorf("Did not add message to buffer")
	}
}

func TestManager_Endpoint_AddProtocol(t *testing.T) {
	pc := &connect.ProtoComms{
		Manager: connect.NewManagerTesting(t),
	}
	m := NewManager(pc, DefaultManagerFlags())
	_, err := m.Endpoint(context.Background(), &GossipMsg{
		Tag:       "test",
		Origin:    []byte("origin"),
		Payload:   []byte("payload"),
		Signature: []byte("signature"),
	})
	if err != nil {
		t.Errorf("Failed to send message: %+v", err)
	}
	record, ok := m.buffer["test"]
	if !ok {
		t.Error("Did not create expected message record")
	}
	if len(record.Messages) != 1 {
		t.Errorf("Did not add message to buffer")
	}

	var received bool
	r := func(msg *GossipMsg) error {
		received = true
		return nil
	}
	v := func(msg *GossipMsg, smth []byte) error {
		return nil
	}
	m.NewGossip("test", DefaultProtocolFlags(), r, v, []*id.ID{})

	if len(m.protocols) != 1 {
		t.Errorf("Failed to add protocol")
	}

	if !received {
		t.Errorf("Did not receive message in buffer")
	}

}

// func TestComms_Stream(t *testing.T) {
//	// TODO: Implement test once streaming is enabled
// }
