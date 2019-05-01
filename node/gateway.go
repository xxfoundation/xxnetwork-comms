////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> gateway functionality

package node

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// SendReceiveBatch sends a batch to the gateway
func SendReceiveBatch(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.Batch) error {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
	ctx, cancel := connect.DefaultContext()

	_, err := c.ReceiveBatch(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ReceiveBatch(): Error received: %+v", err)
	}

	cancel()
	return err
}
