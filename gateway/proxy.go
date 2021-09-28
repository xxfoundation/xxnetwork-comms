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
	"google.golang.org/grpc"
)

// ---------------------- Start of deprecated fields ----------- //

// TODO: Remove comm once RequestClientKey is properly tested
// Gateway -> Gateway forward client RequestNonce.
func (g *Comms) SendRequestNonce(host *connect.Host,
	messages *pb.NonceRequest) (*pb.Nonce, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestNonce(ctx, messages)
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
	result := &pb.Nonce{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// Gateway -> Gateway forward client ConfirmNonce.
// TODO: Remove comm once RequestClientKey is properly tested
func (g *Comms) SendConfirmNonce(host *connect.Host,
	messages *pb.RequestRegistrationConfirmation) (*pb.RegistrationConfirmation, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).ConfirmNonce(ctx, messages)
		if err != nil {
			return nil, err
		}

		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.TRACE.Printf("Sending client ConfirmNonce: %+v", messages)
	resultMsg, err := g.Send(host, f)
	if err != nil {
		return nil, err
	}

	// Marshall the result
	result := &pb.RegistrationConfirmation{}
	return result, ptypes.UnmarshalAny(resultMsg, result)
}

// ---------------------- End of deprecated fields ----------- //

// Gateway -> Gateway forward client RequestNonce.
func (g *Comms) SendRequestClientKey(host *connect.Host,
	messages *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestClientKey(ctx, messages)
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
func (g *Comms) SendPutMessage(host *connect.Host,
	messages *pb.GatewaySlot) (*pb.GatewaySlotResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).PutMessage(ctx, messages)
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
func (g *Comms) SendPutManyMessages(host *connect.Host,
	messages *pb.GatewaySlots) (*pb.GatewaySlotResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).PutManyMessages(ctx, messages)
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
	messages *pb.GetMessages) (*pb.GetMessagesResponse, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := host.GetMessagingContext()
		defer cancel()

		// Send the message
		resultMsg, err := pb.NewGatewayClient(conn).RequestMessages(ctx, messages)
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
