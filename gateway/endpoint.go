////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"golang.org/x/net/context"
)

// Handle a CheckMessages event
func (s *gateway) CheckMessages(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.ClientMessages, error) {
	msgIds, ok := gatewayHandler.CheckMessages(msg.UserID)
	returnMsg := &pb.ClientMessages{}
	if ok {
		returnMsg.MessageIDs = msgIds
	}
	return returnMsg, nil
}

// Handle a GetMessage event
func (s *gateway) GetMessage(ctx context.Context, msg *pb.ClientPollMessage) (
	*pb.CmixMessage, error) {
	returnMsg, ok := gatewayHandler.GetMessage(msg.UserID, msg.MessageID)
	if !ok {
		// Return an empty message if no results
		returnMsg = &pb.CmixMessage{}
	}
	return returnMsg, nil
}
