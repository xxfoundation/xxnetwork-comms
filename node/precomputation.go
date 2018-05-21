////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// precomputation.go - all the comms client functions for precomputation.
package node

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/comms/connect"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SendPrecompShare(addr string, message *pb.PrecompShareMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.PrecompShare(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompShare: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompShareInit(addr string, message *pb.PrecompShareInitMessage) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompShareInit(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompShareInit: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompShareCompare(addr string,
	message *pb.PrecompShareCompareMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompShareCompare(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompShareCompare: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompShareConfirm(addr string,
	message *pb.PrecompShareConfirmMessage) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompShareConfirm(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompShareConfirm: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompDecrypt(addr string, message *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompDecrypt(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompDecrypt: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompEncrypt(addr string, message *pb.PrecompEncryptMessage) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompEncrypt(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompEncrypt: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompPermute(addr string, message *pb.PrecompPermuteMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompPermute(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompPermute: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendPrecompReveal(addr string, message *pb.PrecompRevealMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.PrecompReveal(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))
	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompReveal: Error received: %s", err)
	}
	cancel()
	return result, err
}
