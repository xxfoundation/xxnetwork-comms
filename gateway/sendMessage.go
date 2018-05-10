////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SendGetMessage(addr string, message *pb.ClientPollMessage) (*pb.
	CmixMessage, error) {
	// Attempt to connect to addr
	c := Connect(addr)
	ctx, cancel := DefaultContext()

	// Send the message
	result, err := c.GetMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("GetMessage: Error received: %s", err)
	}
	cancel()
	return result, err
}
