////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> all servers functionality

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Server -> Server Send Function
func (s *Comms) SendGetMeasure(host *connect.Host,
	message *pb.RoundInfo) (*pb.RoundMetrics, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).GetMeasure(ctx, message)
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
	result := &pb.RoundMetrics{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) SendAskOnline(host *connect.Host,
	message *pb.Ping) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).AskOnline(ctx, message)
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) SendFinishRealtime(host *connect.Host,
	message *pb.RoundInfo) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).FinishRealtime(ctx, message)
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) SendNewRound(host *connect.Host,
	message *pb.RoundInfo) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).CreateNewRound(ctx, message)
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) SendPostRoundPublicKey(host *connect.Host,
	message *pb.RoundPublicKey) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).PostRoundPublicKey(ctx, message)
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) SendPostPrecompResult(host *connect.Host,
	roundID uint64, slots []*pb.Slot) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).PostPrecompResult(ctx,
			&pb.Batch{
				Round: &pb.RoundInfo{
					ID: roundID,
				},
				Slots: slots,
			})
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Server -> Server Send Function
func (s *Comms) RoundTripPing(host *connect.Host,
	roundID uint64, payload *any.Any) (*pb.Ack, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).SendRoundTripPing(ctx,
			&pb.RoundTripPing{
				Round: &pb.RoundInfo{
					ID: roundID,
				},
				Payload: payload,
			})
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
	result := &pb.Ack{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
