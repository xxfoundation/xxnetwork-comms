////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	"errors"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Permissioning -> Server Send Function
func (r *RegistrationComms) SendNodeTopology(connInfo *connect.Host,
	message *pb.NodeTopology) error {

	// Obtain the connection
	conn, err := r.ObtainConnection(connInfo)
	if err != nil {
		return err
	}

	// Set up the context
	ctx, cancel := connect.MessagingContext()
	defer cancel()

	// Wrap message as a generic
	anyMessage, err := ptypes.MarshalAny(message)
	if err != nil {
		return err
	}

	// Sign message
	signedMessage, err := r.Manager.SignMessage(anyMessage, "Permissioning")
	if err != nil {
		return err
	}

	// Send the message
	_, err = pb.NewNodeClient(conn.Connection).DownloadTopology(ctx, signedMessage)
	if err != nil {
		err = errors.New(err.Error())
	}

	return err
}
