////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
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
	jww.DEBUG.Printf("Sending Node Registration message: %+v", message)
	_, err := s.Send(host, f)
	return err
}

// Server -> Registration Send Function
func (s *Comms) RequestNdf(host *connect.Host, message *pb.NDFHash) (*pb.NDF, error) {
	// Call Protocomms Request NDF
	return s.ProtoComms.RequestNdf(host, message)

}
