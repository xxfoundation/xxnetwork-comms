///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains gateway -> server registration functionality

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/xx_network/comms/connect"
	"google.golang.org/grpc"
)

// Gateway -> Server Send Function
func (g *Comms) SendRequestNonceMessage(host *connect.Host,
	message *pb.NonceRequest) (*pb.Nonce, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := g.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).RequestNonce(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Request Nonce message: %+v", message)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.Nonce{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Server Send Function
func (g *Comms) SendConfirmNonceMessage(host *connect.Host,
	message *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := g.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).ConfirmRegistration(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Confirm Nonce message: %+v", message)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegistrationConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Server Send Function
func (g *Comms) SendPoll(host *connect.Host,
	message *pb.ServerPoll) (*pb.ServerPollResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()
		//Pack the message for server
		authMsg, err := g.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// Send the message
		resultMsg, err := pb.NewNodeClient(conn).Poll(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending Poll message...")
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.ServerPollResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
