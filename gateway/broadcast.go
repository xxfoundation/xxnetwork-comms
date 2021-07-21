///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains logic for miscellaneous send functions

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"google.golang.org/grpc"
)

// StartDownloadMixedBatch sends a request for streaming a completed batch.
func (g *Comms) StartDownloadMixedBatch(host *connect.Host, ready *pb.BatchReady) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := g.PackAuthenticatedMessage(ready, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).StartDownloadMixedBatch(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message...")
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return err
	}

	// Marshall the result
	result := &messages.Ack{}
	return ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Gateway message sharing within a team
func (g *Comms) SendShareMessages(host *connect.Host, messages *pb.RoundMessages) error {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		authMsg, err := g.PackAuthenticatedMessage(messages, host, false)
		if err != nil {
			return nil, err
		}
		// Send the message
		_, err = pb.NewGatewayClient(conn).ShareMessages(ctx, authMsg)
		if err != nil {
			err = errors.New(err.Error())
		}
		return nil, err
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Share Messages message: %+v", messages)
	_, err := g.Send(host, f)
	return err
}

// GetRoundBufferInfo Asks the server for round buffer info, specifically how
// many rounds have gone through precomputation.
// Note that this function should block if the buffer size is 0
// This allows the caller to continuously poll without spinning too much.
func (g *Comms) GetRoundBufferInfo(host *connect.Host) (*pb.RoundBufferInfo, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack message into an authenticated message
		authMsg, err := g.PackAuthenticatedMessage(&messages.Ping{}, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).GetRoundBufferInfo(ctx,
			authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Get Round Buffer info message...")
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RoundBufferInfo{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
