////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"github.com/pkg/errors"
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
		err = errors.New(err.Error())
		jww.ERROR.Printf("PutMessage: Error received: %+v", err)
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
		err = errors.New(err.Error())
		jww.ERROR.Printf("CheckMessages: Error received: %+v", err)
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
		err = errors.New(err.Error())
		jww.ERROR.Printf("GetMessage: Error received: %+v", err)
	}

	cancel()
	return result, err
}

// Send a RequestNonceMessage to the gateway
func SendRequestNonceMessage(addr string, message *pb.RequestNonceMessage) (
	*pb.NonceMessage, error) {

	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.RequestNonce(ctx, message)

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
func SendConfirmNonceMessage(addr string, message *pb.ConfirmNonceMessage) (
	*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.ConfirmNonce(ctx, message)

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
