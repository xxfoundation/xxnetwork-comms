///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
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

// Gateway -> Gateway message sharing within a team
func (g *Comms) ShareMessages(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {

	authState, err := g.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable to handle reception of AuthenticatedMessage: %+v", err)
	}

	// Marshall the any message to the message type needed
	roundMessages := &pb.RoundMessages{}
	err = ptypes.UnmarshalAny(msg.Message, roundMessages)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, g.handler.ShareMessages(roundMessages, authState)
}

// Ping gateway to ask for users to notify
func (g *Comms) PollForNotifications(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.UserIdList, error) {

	authState, err := g.AuthenticatedReceiver(msg)
	if err != nil {
		return nil, errors.Errorf("Unable to handle reception of AuthenticatedMessage: %+v", err)
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

// Client -> Gateway historical round request
func (g *Comms) RequestHistoricalRounds(ctx context.Context, msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	return g.handler.RequestHistoricalRounds(msg)
}

// Client -> Gateway message request
func (g *Comms) RequestMessages(ctx context.Context, msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	return g.handler.RequestMessages(msg)
}

// Client -> Gateway bloom filter request
func (g *Comms) RequestBloom(ctx context.Context, msg *pb.GetBloom) (*pb.GetBloomResponse, error) {
	return g.handler.RequestBloom(msg)
}
