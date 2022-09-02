////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Endpoints for the gossip protocol manager

package gossip

import (
	"context"
	jww "github.com/spf13/jwalterweatherman"
	"time"
)

// Generic endpoint for forwarding GossipMsg to correct Protocol
func (m *Manager) Endpoint(ctx context.Context, msg *GossipMsg) (*Ack, error) {
	m.protocolLock.RLock()
	defer m.protocolLock.RUnlock()

	if protocol, ok := m.protocols[msg.Tag]; ok {
		// Sometimes the callbacks can block on waiting for appropriate state.
		// This ensures that they will not interfere with the comms endpoint and
		// operate on their own schedule.
		go func(protocol *Protocol, msg *GossipMsg) {
			err := protocol.receive(msg)
			if err != nil {
				jww.ERROR.Printf("Reception of protocol %s encountered an "+
					"error: %+v", msg.Tag, err)
				return
			}
		}(protocol, msg)
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
