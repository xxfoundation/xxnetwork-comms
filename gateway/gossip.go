///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package gateway

import (
	"github.com/pkg/errors"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
)

// Sets up a new gossip manager connected to the gateway's comms object.
// It implicitly initializes a gossip for rate limiting purposes
func (g *Comms) SetUpGossip(gossipFlags gossip.ManagerFlags, receiver gossip.Receiver,
	verifier gossip.SignatureVerification,	peers []*id.ID)  {

	g.Manager = gossip.NewManager(g.ProtoComms, gossipFlags)

	g.Manager.NewGossip("batch", gossip.DefaultProtocolFlags(),
		receiver, verifier, peers)
}

// Add newPeer to the list of gossip peers for the given protocol.
// The protocol is looked up by the given tag. If this protocol is not
// found then an error is returned
func (g *Comms) AddGossipPeer(tag string, newPeer *id.ID) error  {
	p, ok := g.Manager.Get(tag)
	if !ok {
		return errors.Errorf("Could not find gossip protocol with " +
			"tag %s", tag)
	}

	return p.AddGossipPeer(newPeer)
}