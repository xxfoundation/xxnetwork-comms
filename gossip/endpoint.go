///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Endpoints for the gossip protocol manager

package gossip

import (
	"context"
	"time"
)

// Generic endpoint for forwarding GossipMsg to correct Protocol
func (m *Manager) Endpoint(ctx context.Context, msg *GossipMsg) (*Ack, error) {
	m.protocolLock.RLock()
	if protocol, ok := m.protocols[msg.Tag]; ok {
		err := protocol.receive(msg)
		m.protocolLock.RUnlock()
		if err != nil {
			return nil, err
		}
		return &Ack{}, nil
	}

	m.bufferLock.Lock()
	if record, ok := m.buffer[msg.Tag]; ok {
		record.Messages = append(record.Messages, msg)
	} else {
		t := time.Now()
		m.buffer[msg.Tag] = &MessageRecord{
			Timestamp: t,
			Messages:  []*GossipMsg{msg},
		}
	}
	m.bufferLock.Unlock()

	return &Ack{}, nil
}

// Generic streaming endpoint for forwarding GossipMsg to correct Protocol
func (m *Manager) Stream(stream Gossip_StreamServer) error {
	// TODO: Will be implemented later on
	return nil
}
