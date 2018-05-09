////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
)

func RequestContactList(addr string, message *pb.ContactPoll) (*pb.
	ContactMessage, error) {
	c := Connect(addr)
	ctx, cancel := DefaultContext()
	result, err := c.RequestContactList(ctx, message)
	cancel()
	return result, err
}
