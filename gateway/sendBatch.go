////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server functionality

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
func (g *Comms) PostNewBatch(host *connect.Host, messages *pb.Batch) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		_, err := pb.NewNodeClient(conn).PostNewBatch(ctx, messages)
		if err != nil {
			err = errors.New(err.Error())
		}
		return nil, err
	}

	// Execute the Send function
	_, err := host.Send(f)
	return err
}

// GetRoundBufferInfo Asks the server for round buffer info, specifically how
// many rounds have gone through precomputation.
// Note that this function should block if the buffer size is 0
// This allows the caller to continuously poll without spinning too much.
func (g *Comms) GetRoundBufferInfo(host *connect.Host) (*pb.RoundBufferInfo, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).GetRoundBufferInfo(ctx,
			&pb.RoundBufferInfo{})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := host.Send(f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RoundBufferInfo{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Server Send Function
func (g *Comms) GetCompletedBatch(host *connect.Host) (*pb.Batch, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).GetCompletedBatch(ctx,
			&pb.Ping{})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	resultMsg, err := host.Send(f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Batch{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
