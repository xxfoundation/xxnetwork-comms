////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Client -> Gateway Send Function
func (c *Comms) SendPutMessage(connInfo *connect.Host,
	message *pb.Slot) error {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	_, err = pb.NewGatewayClient(conn.Connection).PutMessage(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return err
}

// Client -> Gateway Send Function
func (c *Comms) SendCheckMessages(connInfo *connect.Host,
	message *pb.ClientRequest) (*pb.IDList, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewGatewayClient(conn.Connection).CheckMessages(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Client -> Gateway Send Function
func (c *Comms) SendGetMessage(connInfo *connect.Host,
	message *pb.ClientRequest) (*pb.Slot, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewGatewayClient(conn.Connection).GetMessage(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Client -> Gateway Send Function
func (c *Comms) SendRequestNonceMessage(connInfo *connect.Host,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewGatewayClient(conn.Connection).RequestNonce(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	errMsg := response.GetError()
	if errMsg != "" {
		err = errors.New(errMsg)
	}

	return response, err
}

// Client -> Gateway Send Function
func (c *Comms) SendConfirmNonceMessage(connInfo *connect.Host,
	message *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewGatewayClient(conn.Connection).ConfirmNonce(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}
	errMsg := response.GetError()
	if errMsg != "" {
		err = errors.New(errMsg)
	}

	return response, err
}
