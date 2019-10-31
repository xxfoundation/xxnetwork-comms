////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Gateway -> Server Send Function
func (g *GatewayComms) SendRequestNonceMessage(connInfo *connect.Host,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewNodeClient(conn.Connection).RequestNonce(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return response, err
}

// Gateway -> Server Send Function
func (g *GatewayComms) SendConfirmNonceMessage(connInfo *connect.Host,
	message *pb.RequestRegistrationConfirmation) (
	*pb.RegistrationConfirmation, error) {

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewNodeClient(conn.Connection).ConfirmRegistration(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return response, err
}

// Gateway -> Server Send Function
func (g *GatewayComms) PollSignedCerts(connInfo *connect.Host,
	message *pb.Ping) (*pb.SignedCerts, error) {

	// Obtain the connection
	conn, err := g.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewNodeClient(conn.Connection).GetSignedCert(ctx, message)
	if err != nil {
		return nil, err
	}

	return response, nil
}
