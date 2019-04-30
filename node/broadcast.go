////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> all servers functionality

package node

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

func SendServerMetrics(addr string, serverCertPath string,
	message *pb.ServerMetrics) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr, serverCertPath)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.GetServerMetrics(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ServerMetrics: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendRoundtripPing(addr string, serverCertPath string,
	message *pb.TimePing) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr, serverCertPath)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.RoundtripPing(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RoundtripPing: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendAskOnline(addr string, serverCertPath string, message *pb.Ping) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr, serverCertPath)
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

func SendNewRound(addr string, serverCertPath string, message *pb.Batch) (
	*pb.Ack, error) {
	c := connect.ConnectToNode(addr, serverCertPath)
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

func SendPostRoundPublicKey(addr string, serverCertPath string,
	message *pb.RoundPublicKey) (*pb.Ack, error) {
	c := connect.ConnectToNode(addr, serverCertPath)
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

// SendFinishPrecomputation sends the final message and AD precomputations to
// other nodes.
func SendFinishPrecomputation(addr string, serverCertPath string,
	roundID uint64, slots []*pb.Slot) (*pb.Ack, error) {
	c := connect.ConnectToNode(addr, serverCertPath)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.FinishPrecomputation(ctx,
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
		jww.ERROR.Printf("FinishPrecomputation: Error received: %+v",
			err)
	}

	cancel()
	return result, err
}
