///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
)

// Handles validation of reverse-authentication tokens
func (g *Comms) AuthenticateToken(ctx context.Context,
	msg *messages.AuthenticatedMessage) (*messages.Ack, error) {
	return &messages.Ack{}, g.ValidateToken(msg)
}

// Handles reception of reverse-authentication token requests
func (g *Comms) RequestToken(context.Context, *messages.Ping) (*messages.AssignToken, error) {
	token, err := g.GenerateToken()
	return &messages.AssignToken{
		Token: token,
	}, err
}

// Sends new MessageIDs in the buffer to a client
func (g *Comms) CheckMessages(ctx context.Context,
	msg *pb.ClientRequest) (*pb.IDList, error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := id.Unmarshal(msg.UserID)
	if err != nil {
		return nil, err
	}

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

	userID, err := id.Unmarshal(msg.UserID)
	if err != nil {
		return nil, err
	}

	returnMsg, err := g.handler.GetMessage(userID, msg.LastMessageID, addr)
	if err != nil {
		// Return an empty message if no results
		returnMsg = &pb.Slot{}
	}
	return returnMsg, err
}

// Receives a single message from a client
func (g *Comms) PutMessage(ctx context.Context, msg *pb.GatewaySlot) (*pb.GatewaySlotResponse,
	error) {

	// Get peer information from context
	addr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Upload a message to the cMix Gateway at the peer's IP address
	returnMsg, err := g.handler.PutMessage(msg, addr)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
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
func (g *Comms) PollForNotifications(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.UserIdList, error) {

	authState, err := g.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	ids, err := g.handler.PollForNotifications(authState)
	returnMsg := &pb.UserIdList{}
	if err == nil {
		for i, userID := range ids {
			returnMsg.IDs[i] = userID.Marshal()
		}
	}
	return returnMsg, err
}

// Client -> Gateway unified polling
func (g *Comms) Poll(ctx context.Context, msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
	return g.handler.Poll(msg)
}
