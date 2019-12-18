////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains client -> gateway functionality

package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Client -> Gateway Send Function
func (c *Comms) SendPutMessage(host *connect.Host, message *pb.Slot) error {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		_, err := pb.NewGatewayClient(conn).PutMessage(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
		}
		return nil, err
	}

	// Execute the Send function
	_, err := c.Send(host, f)
	return err
}

// Client -> Gateway Send Function
func (c *Comms) SendCheckMessages(host *connect.Host,
	message *pb.ClientRequest) (*pb.IDList, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).CheckMessages(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.IDList{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendGetMessage(host *connect.Host,
	message *pb.ClientRequest) (*pb.Slot, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).GetMessage(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Slot{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendRequestNonceMessage(host *connect.Host,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestNonce(ctx, message)

		// Make sure there are no errors with sending the message
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Nonce{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> Gateway Send Function
func (c *Comms) SendConfirmNonceMessage(host *connect.Host,
	message *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).ConfirmNonce(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegistrationConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
