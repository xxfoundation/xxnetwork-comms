////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> gateway functionality

package node

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// SendReceiveBatch sends a batch to the gateway
func SendReceiveBatch(addr string, gatewayCertPath string,
	gatewayCertString string, message []*pb.Batch) error {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
	ctx, cancel := connect.DefaultContext()

	outputMessages := pb.Output{Messages: message}

	_, err := c.ReceiveBatch(ctx, &outputMessages)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("ReceiveBatch(): Error received: %s", err)
	}
	cancel()
	return err
}
