package gossip

import (
	"context"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
)

// Object used to embed Gossip functionality in higher-level Comms objects
type Comms struct {
	*connect.ProtoComms
	*Manager
}

// Generic endpoint for forwarding GossipMsg to correct Protocol
func (g *Comms) Endpoint(ctx context.Context, msg *GossipMsg) (*messages.Ack, error) {
	return &messages.Ack{}, g.Get(msg.Tag).receive(msg)
}

// Generic streaming endpoint for forwarding GossipMsg to correct Protocol
func (g *Comms) Stream(ctx context.Context, msg *GossipMsg) (*messages.Ack, error) {
	// TODO: Will be implemented later on
	return &messages.Ack{}, nil
}
