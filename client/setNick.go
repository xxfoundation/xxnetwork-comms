////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/privategrity/comms/connect"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SetNick(addr string, message *pb.Contact) (*pb.Ack, error) {
	c := connect.ConnectToNode(addr)
	ctx, cancel := connect.DefaultContext()
	result, err := c.SetNick(ctx, message)
	cancel()
	return result, err
}
