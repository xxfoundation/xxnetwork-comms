////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// precomputation.go - all the comms client functions for precomputation.
package clusterclient

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"
	jww "github.com/spf13/jwalterweatherman"
)

func SendPrecompShare(addr string, message *pb.PrecompShareMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.PrecompShare(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompShare: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}

func SendPrecompDecrypt(addr string, message *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.PrecompDecrypt(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompDecrypt: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}

func SendPrecompEncrypt(addr string, message *pb.PrecompEncryptMessage) (
	*pb.Ack, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.PrecompEncrypt(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompEncrypt: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}

func SendPrecompPermute(addr string, message *pb.PrecompPermuteMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.PrecompPermute(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompPermute: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}

func SendPrecompReveal(addr string, message *pb.PrecompRevealMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", addr)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	// Send the message
	result, err := c.PrecompReveal(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompReveal: Error received: %s", err)
	}
	cancel()
	conn.Close()

	return result, err
}
