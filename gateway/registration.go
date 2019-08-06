////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a RequestNonceMessage to the server
func (g *GatewayComms) SendRequestNonceMessage(id fmt.Stringer,
	message *pb.NonceRequest) (
	*pb.Nonce, error) {

	// Attempt to connect to addr
	c := g.GetNodeConnection(id)
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
func (g *GatewayComms) SendConfirmNonceMessage(id fmt.Stringer,
	message *pb.DSASignature) (
	*pb.RegistrationConfirmation, error) {

	// Attempt to connect to addr
	c := g.GetNodeConnection(id)
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

// SendGetSignedCertMessage gets signed certs from a node.  Accepts a Ping message (empty message)
func (g *GatewayComms) PollSignedCerts(id fmt.Stringer,
	message *pb.Ping) (*pb.SignedCerts, error) {

	c := g.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	response, err := c.GetSignedCert(ctx, message)

	if err != nil {
		return nil, err
	}

	cancel()
	return response, nil
}
