////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a message to the gateway
func SendPutMessage(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.
	CmixBatch) error {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
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
func SendCheckMessages(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.ClientPollMessage) (*pb.
	ClientMessages, error) {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
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
func SendGetMessage(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.ClientPollMessage) (*pb.
	CmixBatch, error) {
	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
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

// Send a RequestNonceMessage to the gateway
func SendRequestNonceMessage(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.RequestNonceMessage) (
	*pb.NonceMessage, error) {

	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.RequestNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		jww.ERROR.Printf("RequestNonceMessage: Error received: %s", err)
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
func SendConfirmNonceMessage(addr string, gatewayCertPath string,
	gatewayCertString string, message *pb.ConfirmNonceMessage) (
	*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	c := connect.ConnectToGateway(addr, gatewayCertPath, gatewayCertString)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.ConfirmNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		jww.ERROR.Printf("ConfirmNonceMessage: Error received: %s", err)
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
