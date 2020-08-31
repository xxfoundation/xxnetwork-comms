///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package gateway

import (
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

