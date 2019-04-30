////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> server functionality for realtime operations

package node

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

func SendRealtimePermute(addr string, message *pb.RealtimePermuteMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.RealtimePermute(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RealtimePermute: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendRealtimeDecrypt(addr string, message *pb.RealtimeDecryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.RealtimeDecrypt(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RealtimeDecrypt: Error received: %+v", err)
	}

	cancel()
	return result, err
}

func SendRealtimeEncrypt(addr string, message *pb.RealtimeEncryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.RealtimeEncrypt(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("RealtimeEncrypt: Error received: %+v", err)
	}

	cancel()
	return result, err
}
