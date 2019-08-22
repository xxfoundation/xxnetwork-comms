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
