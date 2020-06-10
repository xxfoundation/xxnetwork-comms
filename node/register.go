////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package node

import (
	"github.com/golang/protobuf/ptypes"
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
func (s *Comms) SendPoll(host *connect.Host,
	message *pb.PermissioningPoll) (*pb.PermissionPollResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := s.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewRegistrationClient(conn).Poll(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message...")
	resultMsg, err := s.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.PermissionPollResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Registration Send Function
func (s *Comms) SendRegistrationCheck(host *connect.Host,
	message *pb.RegisteredNodeCheck) (*pb.RegisteredNodeConfirmation, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := s.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewRegistrationClient(conn).CheckRegistration(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)

	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Node Registration Check message: %+v", message)
	resultMsg, err := s.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegisteredNodeConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
