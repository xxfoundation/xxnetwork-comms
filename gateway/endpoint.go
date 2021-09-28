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
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// ---------------------- Start of deprecated fields ----------- //
// TODO: Remove comm once RequestClientKey is properly tested
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

// ---------------------- End of deprecated fields ----------- //

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

// Pass-through for Registration Nonce Communication
func (g *Comms) RequestClientKey(ctx context.Context,
	msg *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {

	return g.handler.RequestClientKey(msg)
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

// Client -> Gateway unified polling
func (g *Comms) Poll(msg *pb.GatewayPoll, stream pb.Gateway_PollServer) error {
	// Get response from higher level
	response, err := g.handler.Poll(msg)
	if err != nil {
		return err
	}

	// Split response into streamable chunks
	chunks, err := pb.SplitResponseIntoChunks(response)
	if err != nil {
		return err
	}

	// Send a header informing client-side of the total number of chunks
	metadataMap := map[string]string{
		pb.ChunkHeader: strconv.Itoa(len(chunks)),
	}

	md := metadata.New(metadataMap)
	if err = stream.SendHeader(md); err != nil {
		return errors.Errorf("Failed to send streaming header: %v", err)
	}

	// Stream each chunk individually
	for i, chunk := range chunks {
		err = stream.Send(chunk)
		if err != nil {
			return errors.Errorf("Failed to send chunk (%d/%d) for "+
				"client polling: %v", i, len(chunks), err)
		}
	}

	return nil
}

// Client -> Gateway historical round request
func (g *Comms) RequestHistoricalRounds(ctx context.Context, msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	return g.handler.RequestHistoricalRounds(msg)
}

// Client -> Gateway message request
func (g *Comms) RequestMessages(ctx context.Context, msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	return g.handler.RequestMessages(msg)
}
