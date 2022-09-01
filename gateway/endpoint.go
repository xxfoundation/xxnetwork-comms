////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains gateway gRPC endpoints

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

// Pass-through for Registration Nonce Communication
func (g *Comms) RequestClientKey(ctx context.Context,
	msg *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {

	return g.handler.RequestClientKey(msg)
}

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

	ipAddr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Upload a message to the cMix Gateway
	returnMsg, err := g.handler.PutMessage(msg, ipAddr)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
}

// Upload many messages to the cMix Gateway
func (g *Comms) PutManyMessages(ctx context.Context, msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse,
	error) {

	ipAddr, _, err := connect.GetAddressFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Upload messages to the cMix Gateway
	returnMsg, err := g.handler.PutManyMessages(msgs, ipAddr)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
}

// Receives a single message from a gateway proxy
func (g *Comms) PutMessageProxy(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.GatewaySlotResponse,
	error) {

	// Verify the message authentication
	authState, err := g.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unnmarshall the any message to the message type needed
	slot := &pb.GatewaySlot{}
	err = ptypes.UnmarshalAny(msg.Message, slot)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Upload a message to the cMix Gateway
	returnMsg, err := g.handler.PutMessageProxy(slot, authState)
	if err != nil {
		returnMsg = &pb.GatewaySlotResponse{}

	}
	return returnMsg, err
}

// Upload many messages to the cMix Gateway from a proxy
func (g *Comms) PutManyMessagesProxy(ctx context.Context, msg *messages.AuthenticatedMessage) (*pb.GatewaySlotResponse,
	error) {

	// Verify the message authentication
	authState, err := g.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		return nil, errors.Errorf("Unable handles reception of AuthenticatedMessage: %+v", err)
	}

	// Unnmarshall the any message to the message type needed
	slots := &pb.GatewaySlots{}
	err = ptypes.UnmarshalAny(msg.Message, slots)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Upload messages to the cMix Gateway
	returnMsg, err := g.handler.PutManyMessagesProxy(slots, authState)
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
