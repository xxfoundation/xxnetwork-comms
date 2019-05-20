////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server -> server functionality for precomputation operations

package node

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

func (s *NodeComms) SendPostPhase(id fmt.Stringer,
	message *pb.Batch) (*pb.Ack, error) {
	// Attempt to connect to addr
	c := s.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	result, err := c.PostPhase(ctx, message,
		grpc_retry.WithMax(connect.MAX_RETRIES))

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("PostPhase: Error received: %+v", err)
	}

	cancel()
	return result, err
}
