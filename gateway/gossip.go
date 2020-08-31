///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////
package gateway

import (
	"encoding/json"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
)

// Sets up a new gossip manager connected to the gateway's comms object.
// It implicitly initializes a gossip for rate limiting purposes
func (g *Comms) SetUpGossip(gossipFlags gossip.ManagerFlags, receiver gossip.Receiver,
	verifier gossip.SignatureVerification, peers []*id.ID) {

	g.Manager = gossip.NewManager(g.ProtoComms, gossipFlags)

	g.Manager.NewGossip("batch", gossip.DefaultProtocolFlags(),
		receiver, verifier, peers)
}

// Add newPeer to the list of gossip peers for the given protocol.
// The protocol is looked up by the given tag. If this protocol is not
// found then an error is returned
func (g *Comms) AddGossipPeer(tag string, newPeer *id.ID) error {
	p, ok := g.Manager.Get(tag)
	if !ok {
		return errors.Errorf("Could not find gossip protocol with "+
			"tag %s", tag)
	}

	return p.AddGossipPeer(newPeer)
}

// gossipBatch builds a gossip message containing all of the sender ID's
// within the batch and gossips it to all peers
func (g *Comms) GossipBatch(batch *pb.Batch) error {

	gossipProtocol, ok := g.Manager.Get("batch")
	if !ok {
		return errors.Errorf("Unable to get gossip protocol. " +
			"Sending batch without gossiping...")

	}

	// Collect all of the sender IDs in the batch
	var senderIds []*id.ID
	for i, slot := range batch.Slots {
		sender, err := id.Unmarshal(slot.SenderID)
		if err != nil {
			return errors.Errorf("Could not completely "+
				"gossip for slot %d in round %d: "+
				"Unreadable sender ID: %v: %v",
				i, batch.Round.ID, slot.SenderID, err)
		}

		senderIds = append(senderIds, sender)
	}

	// Marshal the list of senders into a json
	payloadData, err := json.Marshal(senderIds)
	if err != nil {
		return errors.Errorf("Could not form gossip payload: %v", err)
	}

	var receivedId []*id.ID
	err = json.Unmarshal(payloadData, &receivedId)
	if err != nil {
		return errors.Errorf("Could not marshal sender ID's into a payload: %v", err)
	}

	// Build the message
	gossipMsg := &gossip.GossipMsg{
		Tag:     "batch",
		Origin:  g.Id.Bytes(),
		Payload: payloadData,
	}

	// Sign the gateway message
	sig, err := g.SignMessage(gossipMsg)
	if err != nil {
		return errors.Errorf("Could not sign gossip mesage: %v", err)

	}

	gossipMsg.Signature = sig

	// Gossip the message
	_, errs := gossipProtocol.Gossip(gossipMsg)

	// Return any errors up the stack
	if len(errs) != 0 {
		return errors.Errorf("Could not send to peers: %v", errs)

	}

	return nil
}
