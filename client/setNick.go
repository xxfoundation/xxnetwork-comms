////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/connect"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
)

func SetNick(addr string, message *pb.Contact) (*pb.Ack, error) {
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	result, err := c.SetNick(ctx, message, grpc_retry.WithMax(connect.MAX_RETRIES))
	cancel()
	return result, err
}
