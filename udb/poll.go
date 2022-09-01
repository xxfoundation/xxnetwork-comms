////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains send functions used for polling

package udb

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"google.golang.org/grpc"
)

// RequestNdf is used by User Discovery to Request a NDF from permissioning
func (u *Comms) RequestNdf(host *connect.Host) (*pb.NDF, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewRegistrationClient(
			conn).PollNdf(ctx, &pb.NDFHash{Hash: make([]byte, 0)})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Ndf message...")
	resultMsg, err := u.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &pb.NDF{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
