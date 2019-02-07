////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway GRPC endpoints

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
	"gitlab.com/elixxir/primitives/userid"
)

// Sends new MessageIDs in the buffer to a client
func (s *gateway) CheckMessages(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.ClientMessages, error) {
	userID := new(userid.UserID).SetBytes(msg.UserID)
	msgIds, ok := gatewayHandler.CheckMessages(userID, msg.MessageID)
	returnMsg := &pb.ClientMessages{}
	if ok {
		returnMsg.MessageIDs = msgIds
	}
	return returnMsg, nil
}

// Sends a message matching the given parameters to a client
func (s *gateway) GetMessage(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.CmixMessage, error) {
	userID := new(userid.UserID).SetBytes(msg.UserID)
	returnMsg, ok := gatewayHandler.GetMessage(userID, msg.MessageID)
	if !ok {
		// Return an empty message if no results
		returnMsg = &pb.CmixMessage{}
	}
	return returnMsg, nil
}

// Receives a single message from a client
func (s *gateway) PutMessage(ctx context.Context, msg *pb.CmixMessage) (*pb.Ack,
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
