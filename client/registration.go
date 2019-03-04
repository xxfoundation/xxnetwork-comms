////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> registration server functionality

package client

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a RegisterUserMessage to the RegistrationServer
func SendRegistrationMessage(addr string, message *pb.RegisterUserMessage) error {
	// Attempt to connect to addr
	c := connect.ConnectToRegistration(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	_, err := c.RegisterUser(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("RegistrationMessage: Error received: %s", err)
	}
	cancel()
	return err
}
