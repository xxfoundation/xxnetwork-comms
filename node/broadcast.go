////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> all servers functionality

package node

import (
	"github.com/golang/protobuf/ptypes/any"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Server -> Server Send Function
func (s *NodeComms) SendGetMeasure(connInfo *connect.Host,
	message *pb.RoundInfo) (*pb.RoundMetrics, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).GetMeasure(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))

	// Check for errors
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) SendAskOnline(connInfo *connect.Host,
	message *pb.Ping) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).AskOnline(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) SendFinishRealtime(connInfo *connect.Host,
	message *pb.RoundInfo) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).FinishRealtime(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	cancel()
	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) SendNewRound(connInfo *connect.Host,
	message *pb.RoundInfo) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).CreateNewRound(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) SendPostRoundPublicKey(connInfo *connect.Host,
	message *pb.RoundPublicKey) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).PostRoundPublicKey(ctx, message,
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) SendPostPrecompResult(connInfo *connect.Host,
	roundID uint64, slots []*pb.Slot) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).PostPrecompResult(ctx,
		&pb.Batch{
			Round: &pb.RoundInfo{
				ID: roundID,
			},
			Slots: slots,
		},
		grpc_retry.WithMax(connect.DefaultMaxRetries))
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}

// Server -> Server Send Function
func (s *NodeComms) RoundTripPing(connInfo *connect.Host,
	roundID uint64, payload *any.Any) (*pb.Ack, error) {

	// Obtain the connection
	conn, err := s.ObtainConnection(connInfo)
	if err != nil {
		return nil, err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Send the message
	result, err := pb.NewNodeClient(conn.Connection).SendRoundTripPing(ctx,
		&pb.RoundTripPing{
			Round: &pb.RoundInfo{
				ID: roundID,
			},
			Payload: payload,
		})
	if err != nil {
		err = errors.New(err.Error())
	}

	return result, err
}
