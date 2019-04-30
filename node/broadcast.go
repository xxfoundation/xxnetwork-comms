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
	"golang.org/x/net/context"
)

func SetPublicKey(addr string, message *pb.PublicKeyMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.SetPublicKey(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SetPublicKey: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendServerMetrics(addr string, message *pb.ServerMetricsMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.ServerMetrics(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("ServerMetrics: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendRoundtripPing(addr string, message *pb.TimePing) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
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

func SendAskOnline(addr string, message *pb.Ping) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.AskOnline(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("AskOnline: Error received: %+v", err)
	}
	cancel()
	return result, err
}

func SendNewRound(addr string, message *pb.InitRound) (*pb.Ack, error) {
	c := connect.ConnectToNode(addr)

	// Send the message
	result, err := c.NewRound(context.Background(), message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("NewRound: Error received: %+v", err)
	}
	return result, err
}
