///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains gateway -> gateway proxying functionality

package gateway

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"time"
)

// Gateway -> Gateway forward client RequestClientKey.
func (g *Comms) SendRequestClientKey(host *connect.Host,
	messages *pb.SignedClientKeyRequest, timeout time.Duration) (
	*pb.SignedKeyResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn.GetGrpcConn()).
			RequestClientKey(ctx, messages)
		if err != nil {
			return nil, err
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending client RequestNonce: %+v", messages)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.SignedKeyResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Gateway forward client PutMessage.
func (g *Comms) SendPutMessageProxy(host *connect.Host, messages *pb.GatewaySlot, timeout time.Duration) (*pb.GatewaySlotResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// if messages != nil {
		// 	messages.IpAddr = ipAddr
		// }

		// Pack data into authenticated message
		authMsg, err := g.PackAuthenticatedMessage(messages, host, false)
		if err != nil {
			return nil, err
		}

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn.GetGrpcConn()).
			PutMessageProxy(ctx, authMsg)
		if err != nil {
			return nil, err
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending client PutMessage: %+v", messages)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GatewaySlotResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Gateway forward client PutManyMessages.
func (g *Comms) SendPutManyMessagesProxy(host *connect.Host, messages *pb.GatewaySlots, timeout time.Duration) (*pb.GatewaySlotResponse, error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// if messages != nil {
		// 	messages.IpAddr = ipAddr
		// }

		// Pack data into authenticated message
		authMsg, err := g.PackAuthenticatedMessage(messages, host, false)
		if err != nil {
			return nil, err
		}

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn.GetGrpcConn()).
			PutManyMessagesProxy(ctx, authMsg)
		if err != nil {
			return nil, err
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending client PutMessage: %+v", messages)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GatewaySlotResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Gateway forward client RequestMessages.
func (g *Comms) SendRequestMessages(host *connect.Host,
	messages *pb.GetMessages, timeout time.Duration) (*pb.GetMessagesResponse,
	error) {

	// Create the Send Function
	f := func(conn connect.Connection) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContextWithTimeout(timeout)
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn.GetGrpcConn()).
			RequestMessages(ctx, messages)
		if err != nil {
			return nil, err
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending client RequestMessages: %+v", messages)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.GetMessagesResponse{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}
