////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
)

// Permissioning -> Server Send Function
func (r *Comms) SendNodeTopology(host *connect.Host,
	message *pb.NodeTopology) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := connect.MessagingContext()
		defer cancel()

		// Wrap message as a generic
		anyMessage, err := ptypes.MarshalAny(message)
		if err != nil {
			return nil, err
		}

		// Sign message
		signedMessage, err := r.Manager.SignMessage(anyMessage, "Permissioning")
		if err != nil {
			return nil, err
		}

		// Send the message
		_, err = pb.NewNodeClient(conn).DownloadTopology(ctx, signedMessage)
		if err != nil {
			err = errors.New(err.Error())
		}
		return nil, err
	}

	// Execute the Send function
	_, err := host.Send(f)
	return err
}
