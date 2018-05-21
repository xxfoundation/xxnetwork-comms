////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/comms/connect"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

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

func SendReceiveBatch(addr string, message []*pb.CmixMessage) error {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	outputMessages := pb.OutputMessages{Messages: message}

	_, err := c.ReceiveBatch(ctx, &outputMessages)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("ReceiveBatch(): Error received: %s", err)
	}
	cancel()
	return err
}
