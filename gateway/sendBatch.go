////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server functionality

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Sends a batch of messages from the gateway to a server
func SendBatch(addr string, messages []*pb.CmixMessage) error {
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Create an InputMessage
	msgs := &pb.InputMessages{Messages: messages}

	_, err := c.StartRound(ctx, msgs)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("SendBatch: Error received: %s", err)
	}
	cancel()
	return err
}
