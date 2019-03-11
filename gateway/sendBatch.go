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

// GetRoundBufferInfo Asks the server for round buffer info, specifically how
// many rounds have gone through precomputation.
// Note that this function should block if the buffer size is 0
// This allows the caller to continuously poll without spinning too much.
func GetRoundBufferInfo(addr string) (int, error) {
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	msg := &pb.Ping{}
	bufInfo, err := c.GetRoundBufferInfo(ctx, msg)
	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("GetRoundBufferInfo: Error received: %s", err)
	}
	cancel()

	bufSize := int(bufInfo.RoundBufferSize)
	return bufSize, err
}
