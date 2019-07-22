////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> all servers functionality

package node

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

func (s *NodeComms) SendGetMeasure(id fmt.Stringer, message *pb.RoundInfo) (*pb.RoundMetrics, error) {
	// Attempt to connect to addr
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.GetMeasure(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Check for errors
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("GetMeasure: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func (s *NodeComms) SendAskOnline(id fmt.Stringer, message *pb.Ping) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.AskOnline(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("AskOnline: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func (s *NodeComms) SendFinishRealtime(id fmt.Stringer,
	message *pb.RoundInfo) (*pb.Ack, error) {
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.FinishRealtime(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("FinishRealtime: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func (s *NodeComms) SendNewRound(id fmt.Stringer, message *pb.RoundInfo) (
	*pb.Ack, error) {
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.CreateNewRound(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("NewRound: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func (s *NodeComms) SendPostRoundPublicKey(id fmt.Stringer,
	message *pb.RoundPublicKey) (*pb.Ack, error) {
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.PostRoundPublicKey(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SendPostRoundPublicKey: Error received: %+v", err)
	}

	cancel()
	return result, err
}

// SendPostPrecompResult sends the final message and AD precomputations to
// other nodes.
func (s *NodeComms) SendPostPrecompResult(id fmt.Stringer,
	roundID uint64, slots []*pb.Slot) (*pb.Ack, error) {
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.PostPrecompResult(ctx,
		&pb.Batch{
			Round: &pb.RoundInfo{
				ID: roundID,
			},
			Slots: slots,
		},
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("PostPrecompResult: Error received: %+v",
			err)
	}

	cancel()
	return result, err
}
