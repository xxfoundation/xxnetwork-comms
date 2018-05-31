////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
)

// Handler implementation for the Gateway
type Handler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userID uint64, messageID string) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage(userID uint64, msgID string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage(*pb.CmixMessage) bool
	// ReceiveBatch receives message from a cMix node
	ReceiveBatch(messages *pb.OutputMessages)
}
