////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// SendBatch kicks off a realtime round by sending a batchsize of messages
// at the node
package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/comms/connect"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SendBatch(addr string, messages []*pb.CmixMessage) (error) {
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
