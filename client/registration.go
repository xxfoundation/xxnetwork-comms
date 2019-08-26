////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> registration server functionality

package client

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a RegisterUserMessage to the RegistrationServer
func (c *ClientComms) SendRegistrationMessage(id fmt.Stringer,
	message *pb.UserRegistration) (*pb.UserRegistrationConfirmation, error) {
	// Attempt to connect to addr
	connection := c.GetRegistrationConnection(id)
	ctx, cancel := connect.MessagingContext()

	// Send the message
	response, err := connection.RegisterUser(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RegistrationMessage: Error received: %+v", err)
	}

	cancel()
	return response, err
}

// Call CheckClientVersion on the registration server
func (c *ClientComms) SendCheckClientVersionMessage(id fmt.Stringer,
	message *pb.ClientVersion) (*pb.ClientVersionConfirmation, error) {
	// Get the connection
	connection := c.GetRegistrationConnection(id)
	ctx, cancel := connect.MessagingContext()

	// Send the message
	response, err := connection.CheckClientVersion(ctx, message)

	// Log if we got an error
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("CheckClientVersion: Error received: %+v", err)
	}

	// Finish up
	cancel()
	return response, err
}
