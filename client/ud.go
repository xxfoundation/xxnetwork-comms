////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
)

// Client -> User Discovery Register User Function
func (c *Comms) SendRegisterUser(host *connect.Host, message *pb.UDBUserRegistration) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RegisterUser(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Register Fact Function
func (c *Comms) SendRegisterFact(host *connect.Host, message *pb.FactRegisterRequest) (*pb.FactRegisterResponse, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RegisterFact(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Register Fact message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.FactRegisterResponse{}

	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendConfirmFact(host *connect.Host, message *pb.FactConfirmRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).ConfirmFact(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendRemoveFact(host *connect.Host, message *pb.FactRemovalRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RemoveFact(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete Fact Message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}

// Client -> User Discovery Delete Fact Function
func (c *Comms) SendRemoveUser(host *connect.Host, message *pb.FactRemovalRequest) (*messages.Ack, error) {
	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewUDBClient(conn).RemoveUser(ctx, message)
		if err != nil {
			err = errors.New(err.Error())
			return nil, errors.New(err.Error())

		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Delete Fact Message: %+v", message)
	_, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	return &messages.Ack{}, nil
}
