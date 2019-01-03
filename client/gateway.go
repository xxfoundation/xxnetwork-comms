////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a message to the gateway
func SendPutMessage(addr string, message *pb.CmixMessage) error {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	_, err := c.PutMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PutMessage: Error received: %s", err)
	}
	cancel()
	return err
}

// Request MessageIDs of new messages in the buffer from the gateway
func SendCheckMessages(addr string, message *pb.ClientPollMessage) (*pb.
	ClientMessages, error) {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.CheckMessages(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("CheckMessages: Error received: %s", err)
	}
	cancel()
	return result, err
}

// Request a message with a specific ID from the gateway
func SendGetMessage(addr string, message *pb.ClientPollMessage) (*pb.
	CmixMessage, error) {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.GetMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("GetMessage: Error received: %s", err)
	}
	cancel()
	return result, err
}
