////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package mixclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func SetNick(addr string, message *pb.Contact) (*pb.Ack, error) {
	c := Connect(addr)
	ctx, cancel := DefaultContext()
	result, err := c.SetNick(ctx, message)
	cancel()
	return result, err
}
