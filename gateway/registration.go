////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Gateway -> Server Send Function
func (g *Comms) SendRequestNonceMessage(host *connect.Host,
	message *pb.AuthenticatedMessage) (*pb.Nonce, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).RequestNonce(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Nonce{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Server Send Function
func (g *Comms) SendConfirmNonceMessage(host *connect.Host,
	message *pb.AuthenticatedMessage) (
	*pb.RegistrationConfirmation, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).ConfirmRegistration(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegistrationConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Server Send Function
func (g *Comms) PollNdf(host *connect.Host,
	message *pb.AuthenticatedMessage) (*pb.GatewayNdf, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).PollNdf(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GatewayNdf{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
