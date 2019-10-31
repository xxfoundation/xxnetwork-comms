////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Server -> Registration Send Function
func (s *NodeComms) SendNodeRegistration(connInfo *connect.ConnectionInfo,
	message *pb.NodeRegistration) error {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	_, err = pb.NewRegistrationClient(conn.Connection).RegisterNode(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
	}

	return err
}
