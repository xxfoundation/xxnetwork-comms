////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server functionality

package gateway

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Sends a batch of messages from the gateway to a server
func (g *GatewayComms) SendBatch(nodeID fmt.Stringer, messages []*pb.
	Batch) error {
	c := g.ConnectToNode(nodeID, nil)
	ctx, cancel := connect.DefaultContext()

	// Create an InputMessage
	msgs := &pb.Input{Messages: messages}

	_, err := c.StartRealtime(ctx, msgs)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SendBatch: Error received: %+v", err)
	}

	cancel()
	return err
}

// GetRoundBufferInfo Asks the server for round buffer info, specifically how
// many rounds have gone through precomputation.
// Note that this function should block if the buffer size is 0
// This allows the caller to continuously poll without spinning too much.
func (g *GatewayComms) GetRoundBufferInfo(nodeID fmt.Stringer) (int,
	error) {
	c := g.ConnectToNode(nodeID, nil)
	ctx, cancel := connect.DefaultContext()

	msg := &pb.RoundBufferInfo{}
	bufSize := int(0)
	bufInfo, err := c.GetRoundBufferInfo(ctx, msg)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("GetRoundBufferInfo: Error received: %+v", err)
	} else {
		bufSize = int(bufInfo.RoundBufferSize)
	}

	cancel()
	return bufSize, err
}
