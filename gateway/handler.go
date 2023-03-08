////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Contains the top-level initialization for the Gateway comms API.
// This includes listening on ports, registering GRPC endpoints, and defining
// the associated callback interface.

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"runtime/debug"
)

// Comms object bundles low-level connect.ProtoComms,
// the gossip.Manager protocol, and the endpoint Handler interface.
type Comms struct {
	*gossip.Manager
	*connect.ProtoComms
	handler Handler
	*pb.UnimplementedGatewayServer
	*messages.UnimplementedGenericServer
}

// Handler describes the endpoint callbacks for Gateway.
type Handler interface {
	PutMessage(message *pb.GatewaySlot, ipAddr string) (*pb.GatewaySlotResponse, error)
	PutManyMessages(msgs *pb.GatewaySlots, ipAddr string) (*pb.GatewaySlotResponse, error)
	PutMessageProxy(message *pb.GatewaySlot, auth *connect.Auth) (*pb.GatewaySlotResponse, error)
	PutManyMessagesProxy(msgs *pb.GatewaySlots, auth *connect.Auth) (*pb.GatewaySlotResponse, error)
	Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)
	RequestClientKey(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error)
	RequestTlsCert(message *pb.RequestGatewayCert) (*pb.GatewayCertificate, error)
	BatchNodeRegistration(msg *pb.SignedClientBatchKeyRequest) (*pb.SignedBatchKeyResponse, error)
	RequestBatchMessages(msg *pb.GetMessagesBatch) (*pb.GetMessagesResponseBatch, error)
}

// StartGateway starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection.
func StartGateway(id *id.ID, localServer string, handler Handler,
	certPem, keyPem []byte, gossipFlags gossip.ManagerFlags) *Comms {

	// Initialize the low-level comms listeners
	pc, err := connect.StartCommServer(id, localServer,
		certPem, keyPem, nil)
	if err != nil {
		jww.FATAL.Panicf("Unable to StartCommServer: %+v", err)
	}
	gatewayServer := Comms{
		handler:    handler,
		ProtoComms: pc,
		Manager:    gossip.NewManager(pc, gossipFlags),
	}

	// Register the high-level comms endpoint functionality
	grpcServer := gatewayServer.GetServer()
	pb.RegisterGatewayServer(grpcServer, &gatewayServer)
	messages.RegisterGenericServer(grpcServer, &gatewayServer)
	gossip.RegisterGossipServer(grpcServer, gatewayServer.Manager)

	pc.ServeWithWeb()
	return &gatewayServer
}

// RestartGateway shuts down &restarts the underlying protocomms server,
// re-registers grpc handlers & starts basic listeners again.  Intended for use
// before replacing https certificates
func (g *Comms) RestartGateway() error {
	g.ProtoComms.Shutdown()
	err := g.ProtoComms.Restart()
	if err != nil {
		return err
	}
	// Register the high-level comms endpoint functionality
	grpcServer := g.GetServer()
	pb.RegisterGatewayServer(grpcServer, g)
	messages.RegisterGenericServer(grpcServer, g)
	gossip.RegisterGossipServer(grpcServer, g.Manager)

	g.ProtoComms.ServeWithWeb()
	return nil
}

// implementationFunctions for the Handler interface.
type implementationFunctions struct {
	PutMessage              func(message *pb.GatewaySlot, ipAddr string) (*pb.GatewaySlotResponse, error)
	PutManyMessages         func(msgs *pb.GatewaySlots, ipAddr string) (*pb.GatewaySlotResponse, error)
	Poll                    func(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	RequestHistoricalRounds func(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	RequestMessages         func(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)
	RequestClientKey        func(message *pb.SignedClientKeyRequest) (*pb.SignedKeyResponse, error)
	PutMessageProxy         func(message *pb.GatewaySlot, auth *connect.Auth) (*pb.GatewaySlotResponse, error)
	PutManyMessagesProxy    func(msgs *pb.GatewaySlots, auth *connect.Auth) (*pb.GatewaySlotResponse, error)
	RequestTlsCert          func(message *pb.RequestGatewayCert) (*pb.GatewayCertificate, error)
	BatchNodeRegistration   func(msg *pb.SignedClientBatchKeyRequest) (*pb.SignedBatchKeyResponse, error)
	RequestBatchMessages    func(msg *pb.GetMessagesBatch) (*pb.GetMessagesResponseBatch, error)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions.
type Implementation struct {
	Functions implementationFunctions
}

// NewImplementation creates and returns a new Handler interface for implementing endpoint callbacks.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			PutMessage: func(message *pb.GatewaySlot, ipAddr string) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return new(pb.GatewaySlotResponse), nil
			},
			PutManyMessages: func(msgs *pb.GatewaySlots, ipAddr string) (*pb.GatewaySlotResponse, error) {
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
			PutMessageProxy: func(message *pb.GatewaySlot, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return &pb.GatewaySlotResponse{}, nil
			},
			PutManyMessagesProxy: func(msgs *pb.GatewaySlots, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return &pb.GatewaySlotResponse{}, nil
			},
			RequestTlsCert: func(message *pb.RequestGatewayCert) (*pb.GatewayCertificate, error) {
				warn(um)
				return &pb.GatewayCertificate{}, nil
			},
			BatchNodeRegistration: func(msg *pb.SignedClientBatchKeyRequest) (*pb.SignedBatchKeyResponse, error) {
				warn(um)
				return &pb.SignedBatchKeyResponse{}, nil
			},
			RequestBatchMessages: func(msg *pb.GetMessagesBatch) (*pb.GetMessagesResponseBatch, error) {
				warn(um)
				return &pb.GetMessagesResponseBatch{}, nil
			},
		},
	}
}

// RequestClientKey is a pass-through for RequestClientKey Communication.
func (s *Implementation) RequestClientKey(message *pb.SignedClientKeyRequest) (
	*pb.SignedKeyResponse, error) {
	return s.Functions.RequestClientKey(message)
}

// PutMessage uploads a message to the cMix Gateway.
func (s *Implementation) PutMessage(message *pb.GatewaySlot, ipAddr string) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutMessage(message, ipAddr)
}

// PutManyMessages uploads many messages to the cMix Gateway.
func (s *Implementation) PutManyMessages(msgs *pb.GatewaySlots, ipAddr string) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutManyMessages(msgs, ipAddr)
}

// PutMessageProxy uploads a message to the cMix Gateway from a proxy gateway.
func (s *Implementation) PutMessageProxy(message *pb.GatewaySlot, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutMessageProxy(message, auth)
}

// PutManyMessagesProxy uploads many messages to the cMix Gateway from a proxy gateway.
func (s *Implementation) PutManyMessagesProxy(msgs *pb.GatewaySlots, auth *connect.Auth) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutManyMessagesProxy(msgs, auth)
}

// Poll provides Client -> Gateway unified polling.
func (s *Implementation) Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error) {
	return s.Functions.Poll(msg)
}

// RequestHistoricalRounds provides Client -> Gateway historical round requests.
func (s *Implementation) RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error) {
	return s.Functions.RequestHistoricalRounds(msg)
}

// RequestMessages handles Client -> Gateway requests for message pickup.
func (s *Implementation) RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error) {
	return s.Functions.RequestMessages(msg)
}

func (s *Implementation) RequestTlsCert(msg *pb.RequestGatewayCert) (*pb.GatewayCertificate, error) {
	return s.Functions.RequestTlsCert(msg)
}

func (s *Implementation) BatchNodeRegistration(msg *pb.SignedClientBatchKeyRequest) (*pb.SignedBatchKeyResponse, error) {
	return s.Functions.BatchNodeRegistration(msg)
}

func (s *Implementation) RequestBatchMessages(msg *pb.GetMessagesBatch) (*pb.GetMessagesResponseBatch, error) {
	return s.Functions.RequestBatchMessages(msg)
}
