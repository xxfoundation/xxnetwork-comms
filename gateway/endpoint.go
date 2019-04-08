////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway GRPC endpoints

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/net/context"
)

// Sends new MessageIDs in the buffer to a client
func (s *gateway) CheckMessages(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.ClientMessages, error) {
	userID := new(id.User).SetBytes(msg.UserID)
	msgIds, ok := gatewayHandler.CheckMessages(userID, msg.MessageID)
	returnMsg := &pb.ClientMessages{}
	if ok {
		returnMsg.MessageIDs = msgIds
	}
	return returnMsg, nil
}

// Sends a message matching the given parameters to a client
func (s *gateway) GetMessage(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.CmixBatch, error) {
	userID := new(id.User).SetBytes(msg.UserID)
	returnMsg, ok := gatewayHandler.GetMessage(userID, msg.MessageID)
	if !ok {
		// Return an empty message if no results
		returnMsg = &pb.CmixBatch{}
	}
	return returnMsg, nil
}

// Receives a single message from a client
func (s *gateway) PutMessage(ctx context.Context, msg *pb.CmixBatch) (*pb.Ack,
	error) {
	gatewayHandler.PutMessage(msg)
	return &pb.Ack{}, nil
}

// Receives a batch of messages from a server
func (s *gateway) ReceiveBatch(ctx context.Context, msg *pb.OutputMessages) (*pb.Ack,
	error) {
	gatewayHandler.ReceiveBatch(msg)
	return &pb.Ack{}, nil
}

// Pass-through for Registration Nonce Communication
func (s *gateway) RequestNonce(ctx context.Context,
	msg *pb.RequestNonceMessage) (*pb.NonceMessage, error) {
	return gatewayHandler.RequestNonce(msg)
}

// Pass-through for Registration Nonce Confirmation
func (s *gateway) ConfirmNonce(ctx context.Context,
	msg *pb.ConfirmNonceMessage) (*pb.RegistrationConfirmation, error) {
	return gatewayHandler.ConfirmNonce(msg)
}
