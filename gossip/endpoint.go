package gossip

import (
	"context"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"time"
)

// Object used to embed Gossip functionality in higher-level Comms objects
type Comms struct {
	*connect.ProtoComms
	*Manager
}

// Generic endpoint for forwarding GossipMsg to correct Protocol
func (g *Comms) Endpoint(ctx context.Context, msg *GossipMsg) (*messages.Ack, error) {

	g.protocolLock.RLock()
	if protocol, ok := g.protocols[msg.Tag]; ok {
		err := protocol.receive(msg)
		g.protocolLock.RUnlock()
		if err != nil {
			return nil, err
		}
		return &messages.Ack{}, nil
	}

	g.bufferLock.Lock()
	if record, ok := g.buffer[msg.Tag]; ok {
		record.Messages = append(record.Messages, msg)
	} else {
		t := time.Now()
		g.buffer[msg.Tag] = &MessageRecord{
			Timestamp: t,
			Messages:  []*GossipMsg{msg},
		}
	}
	g.bufferLock.Unlock()

	return &messages.Ack{}, nil
}

// Generic streaming endpoint for forwarding GossipMsg to correct Protocol
func (g *Comms) Stream(ctx context.Context, msg *GossipMsg) (*messages.Ack, error) {
	// TODO: Will be implemented later on
	return &messages.Ack{}, nil
}
