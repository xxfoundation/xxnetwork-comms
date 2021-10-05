///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains callback interface for gateway functionality

package gateway

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Handler interface for the Gateway
type Handler interface {
	// Upload a message to the cMix Gateway
	PutMessage(message *pb.GatewaySlot) (*pb.GatewaySlotResponse, error)
	// Upload many messages to the cMix Gateway
	PutManyMessages(msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse, error)
	// Client -> Gateway unified polling
	Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	// Client -> Gateway historical round request
	RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	// Client -> Gateway message request
	RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)

	RequestClientKey(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error)
}

// Gateway object used to implement endpoints and top-level comms functionality
type Comms struct {
	*gossip.Manager
	*connect.ProtoComms
	handler Handler
}

// Starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection
func StartGateway(id *id.ID, localServer string, handler Handler,
	certPem, keyPem []byte, gossipFlags gossip.ManagerFlags) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPem, keyPem, nil)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	gatewayServer := Comms{
		handler:    handler,
		ProtoComms: pc,
		Manager:    gossip.NewManager(pc, gossipFlags),
	}

	go func() {
		pb.RegisterGatewayServer(gatewayServer.LocalServer, &gatewayServer)
		messages.RegisterGenericServer(gatewayServer.LocalServer, &gatewayServer)
		gossip.RegisterGossipServer(gatewayServer.LocalServer, gatewayServer.Manager)

		// Register reflection service on gRPC server.
		// This blocks for the lifetime of the listener.
		reflection.Register(gatewayServer.LocalServer)
		if err := gatewayServer.LocalServer.Serve(lis); err != nil {
			jww.FATAL.Panicf("Failed to serve: %+v",
				errors.New(err.Error()))
		}
		jww.INFO.Printf("Shutting down gateway server listener: %s",
			lis)

	}()

	return &gatewayServer
}

// Handler implementation for the Gateway
type implementationFunctions struct {
	// Upload a message to the cMix Gateway
	PutMessage func(message *pb.GatewaySlot) (*pb.GatewaySlotResponse, error)
	// Upload many messages to the cMix Gateway
	PutManyMessages func(msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse, error)
	// Client -> Gateway unified polling
	Poll func(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	// Client -> Gateway historical round request
	RequestHistoricalRounds func(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	// Client -> Gateway message request
	RequestMessages func(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)

	// Pass-through for RequestClientKey Communication
	RequestClientKey func(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// Creates and returns a new Handler interface
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			PutMessage: func(message *pb.GatewaySlot) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return new(pb.GatewaySlotResponse), nil
			},
			PutManyMessages: func(msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return &pb.GatewaySlotResponse{}, nil
			},

			Poll: func(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
				warn(um)
				return &pb.GatewayPollResponse{}, nil
			},
			RequestHistoricalRounds: func(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
				warn(um)
				return &pb.HistoricalRoundsResponse{}, nil
			},
			RequestMessages: func(msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
				warn(um)
				return &pb.GetMessagesResponse{}, nil
			},

			RequestClientKey: func(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error) {
				warn(um)
				return new(pb.SignedKeyResponse), nil
			},
		},
	}
}

// Pass-through for RequestClientKey Communication
func (s *Implementation) RequestClientKey(message *pb.SignedClientKeyRequest) (
	*pb.SignedKeyResponse, error) {
	return s.Functions.RequestClientKey(message)
}

// Upload a message to the cMix Gateway
func (s *Implementation) PutMessage(message *pb.GatewaySlot) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutMessage(message)
}

// Upload many messages to the cMix Gateway
func (s *Implementation) PutManyMessages(msgs *pb.GatewaySlots) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutManyMessages(msgs)
}

// Client -> Gateway unified polling
func (s *Implementation) Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
	return s.Functions.Poll(msg)
}

// Client -> Gateway historical round request
func (s *Implementation) RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	return s.Functions.RequestHistoricalRounds(msg)
}

// Client -> Gateway historical round request
func (s *Implementation) RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	return s.Functions.RequestMessages(msg)
}
