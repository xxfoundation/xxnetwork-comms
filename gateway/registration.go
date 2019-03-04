////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a RequestNonceMessage to the server
func SendRequestNonceMessage(addr string, message *pb.RequestNonceMessage) (
	*pb.NonceMessage, error) {

	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.RequestNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		jww.ERROR.Printf("RequestNonceMessage: Error received: %s", err)
	}

	// Return the NonceMessage
	cancel()
	return response, err
}

// Send a ConfirmNonceMessage to the server
func SendConfirmNonceMessage(addr string, message *pb.ConfirmNonceMessage) (
	*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.ConfirmNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		jww.ERROR.Printf("ConfirmNonceMessage: Error received: %s", err)
	}

	// Return the RegistrationConfirmation
	cancel()
	return response, err
}
