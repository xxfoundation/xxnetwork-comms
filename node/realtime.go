////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// realtime.go - all the comms client functions for realtime.
package node

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/comms/connect"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SendRealtimePermute(addr string, message *pb.RealtimePermuteMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.RealtimePermute(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("RealtimePermute: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendRealtimeDecrypt(addr string, message *pb.RealtimeDecryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.RealtimeDecrypt(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("RealtimeDecrypt: Error received: %s", err)
	}
	cancel()
	return result, err
}

func SendRealtimeEncrypt(addr string, message *pb.RealtimeEncryptMessage) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	// Send the message
	result, err := c.RealtimeEncrypt(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("RealtimeEncrypt: Error received: %s", err)
	}
	cancel()
	return result, err
}
