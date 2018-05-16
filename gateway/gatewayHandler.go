////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
)

type GatewayHandler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userId uint64) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage(userId uint64, msgId string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage(*pb.CmixMessage) bool
}
