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

	// Upload a message to the cMix Gateway
	returnMsg, err := g.handler.PutMessage(msg)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
}

// Upload many messages to the cMix Gateway
func (g *Comms) PutManyMessages(ctx context.Context, msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse,
	error) {

	// Upload messages to the cMix Gateway
	returnMsg, err := g.handler.PutManyMessages(msgs)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
}

// Pass-through for Registration Nonce Communication
func (g *Comms) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {

	return g.handler.RequestNonce(msg)
}

// Pass-through for Registration Nonce Confirmation
func (g *Comms) ConfirmNonce(ctx context.Context,
	msg *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	return g.handler.ConfirmNonce(msg)
}

// Gateway -> Gateway message sharing within a team
func (g *Comms) ShareMessages(ctx context.Context, msg *messages.AuthenticatedMessage) (*messages.Ack, error) {

	authState, err := g.AuthenticatedReceiver(msg, ctx)
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

//// DownloadMixedBatch is the handler for server sending a completed batch to its gateway
//func (g *Comms) DownloadMixedBatch(server pb.Gateway_DownloadMixedBatchServer) error {
//	// Extract the authentication info
//	authMsg, err := connect.UnpackAuthenticatedContext(server.Context())
//	if err != nil {
//		return errors.Errorf("Unable to extract authentication info: %+v", err)
//	}
//
//	authState, err := g.AuthenticatedReceiver(authMsg, server.Context())
//	if err != nil {
//		return errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
//	}
//
//	// Verify the message authentication
//	return g.handler.DownloadMixedBatch(server, authState)
//}
