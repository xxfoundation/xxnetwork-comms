////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a message to the gateway
func (c *ClientComms) SendPutMessage(id fmt.Stringer,
	message *pb.Slot) error {
	// Attempt to connect to addr
	connection := c.GetGatewayConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	_, err := connection.PutMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("PutMessage: Error received: %+v", err)
	}

	cancel()
	return err
}

// Request MessageIDs of new messages in the buffer from the gateway
func (c *ClientComms) SendCheckMessages(id fmt.Stringer,
	message *pb.ClientRequest) (*pb.IDList, error) {
	// Attempt to connect to addr
	connection := c.GetGatewayConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := connection.CheckMessages(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("CheckMessages: Error received: %+v", err)
	}

	cancel()
	return result, err
}

// Request a message with a specific ID from the gateway
func (c *ClientComms) SendGetMessage(id fmt.Stringer,
	message *pb.ClientRequest) (*pb.Slot, error) {
	// Attempt to connect to addr
	connection := c.GetGatewayConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := connection.GetMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("GetMessage: Error received: %+v", err)
	}

	cancel()
	return result, err
}

// Send a RequestNonceMessage to the gateway
func (c *ClientComms) SendRequestNonceMessage(id fmt.Stringer,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Attempt to connect to addr
	connection := c.GetGatewayConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := connection.RequestNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RequestNonceMessage: Error received: %+v", err)
	}

	// Handle logic errors
	errMsg := response.GetError()
	if errMsg != "" {
		jww.ERROR.Printf("RequestNonceMessage: Error received: %s",
			errMsg)
		err = errors.New(errMsg)
	}

	cancel()
	return response, err
}

// Send a ConfirmNonceMessage to the gateway
func (c *ClientComms) SendConfirmNonceMessage(id fmt.Stringer,
	message *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	connection := c.GetGatewayConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message

	response, err := connection.ConfirmNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ConfirmNonceMessage: Error received: %+v", err)
	}

	// Handle logic errors
	errMsg := response.GetError()
	if errMsg != "" {
		jww.ERROR.Printf("ConfirmNonceMessage: Error received: %s",
			errMsg)
		err = errors.New(errMsg)
	}

	cancel()
	return response, err
}
