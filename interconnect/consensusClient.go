///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Package interconnect contains logic for cross communication between the cMix servers and
// xxNetwork's consensus
package interconnect

import (
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// CMixServer -> consensus node Send Function
func (c *CMixServer) GetNdf(host *connect.Host,
	message *messages.Ping) (*NDF, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*anypb.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Format to authenticated message type
		// Send the message

		resultMsg, err := NewInterconnectClient(conn).GetNDF(ctx, message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return anypb.New(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Post Phase message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &NDF{}
	return result, resultMsg.UnmarshalTo(result)
}
