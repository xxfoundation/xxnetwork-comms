////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a RequestNonceMessage to the server
func SendRequestNonceMessage(addr string,
	serverCertPath string, message *pb.NonceRequest) (
	*pb.Nonce, error) {

	// Attempt to connect to addr
	c := connect.ConnectToNode(addr, serverCertPath)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.RequestNonce(ctx, message)

	// Handle comms errors
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RequestNonceMessage: Error received: %+v", err)
	}

	// Return the NonceMessage
	cancel()
	return response, err
}

// Send a ConfirmNonceMessage to the server
func SendConfirmNonceMessage(addr string,
	serverCertPath string, message *pb.DSASignature) (
	*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	c := connect.ConnectToNode(addr, serverCertPath)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	response, err := c.ConfirmRegistration(ctx, message)

	// Handle comms errors
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ConfirmNonceMessage: Error received: %+v", err)
	}

	// Return the RegistrationConfirmation
	cancel()
	return response, err
}
