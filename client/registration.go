////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> registration server functionality

package client

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Client -> Registration Send Function
func (c *ClientComms) SendRegistrationMessage(connInfo *connect.Host,
	message *pb.UserRegistration) (*pb.UserRegistrationConfirmation, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewRegistrationClient(
		conn.Connection).RegisterUser(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return response, err
}

// Client -> Registration Send Function
func (c *ClientComms) SendGetCurrentClientVersionMessage(
	connInfo *connect.Host) (*pb.ClientVersion, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewRegistrationClient(
		conn.Connection).GetCurrentClientVersion(ctx, &pb.Ping{})
	if err != nil {
		err = errors.New(err.Error())
	}

	return response, err
}

// Client -> Registration Send Function
func (c *ClientComms) SendGetUpdatedNDF(connInfo *connect.Host,
	message *pb.NDFHash) (*pb.NDF, error) {

	// Obtain the connection
	conn, err := c.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	response, err := pb.NewRegistrationClient(
		conn.Connection).GetUpdatedNDF(ctx, message)
	if err != nil {
		err = errors.New(err.Error())
	}

	return response, err
}
