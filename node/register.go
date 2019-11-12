////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"errors"
	"github.com/golang/protobuf/ptypes/any"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Server -> Registration Send Function
func (s *Comms) SendNodeRegistration(host *connect.Host,
	message *pb.NodeRegistration) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		_, err := pb.NewRegistrationClient(conn).RegisterNode(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
		}
		return nil, err
	}

	// Execute the Send function
	_, err := host.Send(f)
	return err
}
