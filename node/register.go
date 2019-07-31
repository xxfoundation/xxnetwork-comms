////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"errors"
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a message to the gateway
func (s *NodeComms) SendNodeRegistration(id fmt.Stringer,
	message *pb.NodeRegistration) error {

	// Attempt to connect to addr
	connection := s.GetRegistrationConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	_, err := connection.RegisterNode(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SendNodeRegistration: Error received: %+v", err)
	}

	cancel()
	return err
}
