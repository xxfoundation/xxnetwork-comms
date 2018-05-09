////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Send and receive comms functionality for cMix clients
package client

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

// SendMessageToServer sends a user's message to the cMix cluster
func SendMessageToServer(addr string, message *pb.CmixMessage) (*pb.Ack, error) {
	c := Connect(addr)
	ctx, cancel := DefaultContext()
	result, err := c.ClientSendMessageToServer(ctx, message)
	cancel()

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("SendMessageToServer: Error received: %s", err)
	}
	return result, err
}

// SendClientPoll polls the server for new messages
func SendClientPoll(addr string, message *pb.ClientPollMessage) (*pb.CmixMessage, error) {
	c := Connect(addr)
	ctx, cancel := DefaultContext()
	// Send the message
	result, err := c.ClientPoll(ctx, message)
	cancel()

	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("SendClientPoll: Error received: %s", err)
	}
	return result, err
}
