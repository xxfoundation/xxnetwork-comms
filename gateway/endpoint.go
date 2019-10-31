////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
)

// Sends new MessageIDs in the buffer to a client
func (g *GatewayComms) CheckMessages(ctx context.Context, msg *pb.ClientRequest) (
	*pb.IDList, error) {

	// Get peer information from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.IDList{}, nil
	}

	// Strip port from IP address
	ipAddress, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	userID := id.NewUserFromBytes(msg.UserID)
	msgIds, err := g.handler.CheckMessages(userID, msg.LastMessageID, ipAddress)
	returnMsg := &pb.IDList{}
	if err == nil {
		returnMsg.IDs = msgIds
	}
	return returnMsg, err
}

// Sends a message matching the given parameters to a client
func (g *GatewayComms) GetMessage(ctx context.Context, msg *pb.ClientRequest) (
	*pb.Slot, error) {
	// Get peer information from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Slot{}, nil
	}

	// Strip port from IP address
	ipAddress, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	userID := id.NewUserFromBytes(msg.UserID)
	returnMsg, err := g.handler.GetMessage(userID, msg.LastMessageID, ipAddress)
	if err != nil {
		// Return an empty message if no results
		returnMsg = &pb.Slot{}
	}
	return returnMsg, err
}

// Receives a single message from a client
func (g *GatewayComms) PutMessage(ctx context.Context, msg *pb.Slot) (*pb.Ack,
	error) {

	// Get peer information from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Ack{}, nil
	}

	// Strip port from IP address
	ipAddress, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	// Upload a message to the cMix Gateway at the peer's IP address
	g.handler.PutMessage(msg, ipAddress)

	return &pb.Ack{}, nil
}

// Pass-through for Registration Nonce Communication
func (g *GatewayComms) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {

	// Get peer information from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.Nonce{}, nil
	}

	// Strip port from IP address
	ipAddress, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	return g.handler.RequestNonce(msg, ipAddress)
}

// Pass-through for Registration Nonce Confirmation
func (g *GatewayComms) ConfirmNonce(ctx context.Context,
	msg *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation,
	error) {

	// Get peer information from context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.RegistrationConfirmation{}, nil
	}

	// Strip port from IP address
	ipAddress, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}

	return g.handler.ConfirmNonce(msg, ipAddress)
}
