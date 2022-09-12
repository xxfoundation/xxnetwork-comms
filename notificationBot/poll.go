///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains notificationBot -> all servers functionality

package notificationBot

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
)

// PollNdf gets the NDF from the permissioning server
func (nb *Comms) PollNdf(host *connect.Host, ndfHash []byte) (*pb.NDF, error) {
	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// We use an empty NDF Hash to request an NDF
		ndfRequest := &pb.NDFHash{Hash: ndfHash}

		// Send the message
		resultMsg, err := pb.NewRegistrationClient(conn.GetGrpcConn()).
			PollNdf(ctx, ndfRequest)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Ndf message...")
	resultMsg, err := nb.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.NDF{}
	err = ptypes.UnmarshalAny(resultMsg, result)
	return result, err
}
