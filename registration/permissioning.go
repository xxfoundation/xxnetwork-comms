////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package registration

import (
	"errors"
	"fmt"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Send a message to the gateway
func (r *RegistrationComms) SendNodeTopology(id fmt.Stringer,
	message *pb.NodeTopology) error {

	// Attempt to connect to addr
	connection := r.GetNodeConnection(id)
	ctx, cancel := connect.DefaultContext()

	// Send the message
	_, err := connection.DownloadTopology(ctx, message)

	// Make sure there are no errors with sending the message
	if err != nil {
		err = errors.New(err.Error())
		jww.ERROR.Printf("SendNodeToplogy: Error received: %+v", err)
	}

	cancel()
	return err
}
