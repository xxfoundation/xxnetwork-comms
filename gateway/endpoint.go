////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/net/context"
)

// Sends new MessageIDs in the buffer to a client
func (g *Comms) CheckMessages(ctx context.Context,
	msg *pb.ClientRequest) (*pb.IDList, error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID := id.NewUserFromBytes(msg.UserID)
	msgIds, err := g.handler.CheckMessages(userID, msg.LastMessageID, addr)
	returnMsg := &pb.IDList{}
	if err == nil {
		returnMsg.IDs = msgIds
	}
	return returnMsg, err
}

// Sends a message matching the given parameters to a client
func (g *Comms) GetMessage(ctx context.Context, msg *pb.ClientRequest) (
	*pb.Slot, error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userID := id.NewUserFromBytes(msg.UserID)
	returnMsg, err := g.handler.GetMessage(userID, msg.LastMessageID, addr)
	if err != nil {
		// Return an empty message if no results
		returnMsg = &pb.Slot{}
	}
	return returnMsg, err
}

// Receives a single message from a client
func (g *Comms) PutMessage(ctx context.Context, msg *pb.Slot) (*pb.Ack,
	error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Upload a message to the cMix Gateway at the peer's IP address
	err = g.handler.PutMessage(msg, addr)

	return &pb.Ack{}, err
}

// Pass-through for Registration Nonce Communication
func (g *Comms) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return g.handler.RequestNonce(msg, addr)
}

// Pass-through for Registration Nonce Confirmation
func (g *Comms) ConfirmNonce(ctx context.Context,
	msg *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation,
	error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return g.handler.ConfirmNonce(msg, addr)
}

// Ping gateway to ask for users to notify
func (g *Comms) PollForNotifications(ctx context.Context, msg *pb.Ping) (*pb.IDList, error) {

	ids, err := g.handler.PollForNotifications()
	returnMsg := &pb.IDList{}
	if err == nil {
		returnMsg.IDs = ids
	}
	return returnMsg, err
}
