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
	// Returns the message matching the given parameters to the client
	GetMessage(userId uint64, msgId string) (*pb.CmixMessage, bool)
}
