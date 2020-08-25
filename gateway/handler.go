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
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Handler interface for the Gateway
type Handler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userID *id.ID, messageID string, ipAddress string) ([]string, error)
	// Returns the message matching the given parameters to the client
	GetMessage(userID *id.ID, msgID string, ipAddress string) (*pb.Slot, error) // todo: depracate?
	// Upload a message to the cMix Gateway
	PutMessage(message *pb.GatewaySlot, ipAddress string) (*pb.GatewaySlotResponse, error)
	// Pass-through for Registration Nonce Communication
	RequestNonce(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error)
	// Pass-through for Registration Nonce Confirmation
	ConfirmNonce(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.
		RegistrationConfirmation, error)
	// Ping gateway to ask for users to notify
	PollForNotifications(auth *connect.Auth) ([]*id.ID, error)
	// Client -> Gateway unified polling
	Poll(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	// Client -> Gateway historical round request
	RequestHistoricalRounds(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	// Client -> Gateway message request
	RequestMessages(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)
	// Client -> Gateway bloom request
	RequestBloom(msg *pb.GetBloom) (*pb.GetBloomResponse, error)
}

// Gateway object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

// Starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection
func StartGateway(id *id.ID, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	gatewayServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		pb.RegisterGatewayServer(gatewayServer.LocalServer, &gatewayServer)
		messages.RegisterGenericServer(gatewayServer.LocalServer, &gatewayServer)

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
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages func(userID *id.ID, messageID string, ipAddress string) ([]string, error)
	// Returns the message matching the given parameters to the client
	GetMessage func(userID *id.ID, msgID string, ipAddress string) (*pb.Slot, error)
	// Upload a message to the cMix Gateway
	PutMessage func(message *pb.GatewaySlot, ipAddress string) (*pb.GatewaySlotResponse, error)
	// Pass-through for Registration Nonce Communication
	RequestNonce func(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error)
	// Pass-through for Registration Nonce Confirmation
	ConfirmNonce func(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.
			RegistrationConfirmation, error)
	// Ping gateway to ask for users to notify
	PollForNotifications func(auth *connect.Auth) ([]*id.ID, error)
	// Client -> Gateway unified polling
	Poll func(msg *pb.GatewayPoll) (*pb.GatewayPollResponse, error)
	// Client -> Gateway historical round request
	RequestHistoricalRounds func(msg *pb.HistoricalRounds) (*pb.HistoricalRoundsResponse, error)
	// Client -> Gateway message request
	RequestMessages func(msg *pb.GetMessages) (*pb.GetMessagesResponse, error)
	// Client -> Gateway bloom request
	RequestBloom func(msg *pb.GetBloom) (*pb.GetBloomResponse, error)
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
			CheckMessages: func(userID *id.ID, messageID string, ipAddress string) ([]string, error) {
				warn(um)
				return nil, nil
			},
			GetMessage: func(userID *id.ID, msgID string, ipAddress string) (*pb.Slot, error) {
				warn(um)
				return &pb.Slot{}, nil
			},
			PutMessage: func(message *pb.GatewaySlot, ipAddress string) (*pb.GatewaySlotResponse, error) {
				warn(um)
				return new(pb.GatewaySlotResponse), nil
			},
			RequestNonce: func(message *pb.NonceRequest, ipAddress string) (*pb.Nonce, error) {
				warn(um)
				return new(pb.Nonce), nil
			},
			ConfirmNonce: func(message *pb.RequestRegistrationConfirmation, ipAddress string) (*pb.RegistrationConfirmation, error) {
				warn(um)
				return new(pb.RegistrationConfirmation), nil
			},
			PollForNotifications: func(auth *connect.Auth) ([]*id.ID, error) {
				warn(um)
				return nil, nil
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
			RequestBloom: func(msg *pb.GetBloom) (*pb.GetBloomResponse, error) {
				warn(um)
				return &pb.GetBloomResponse{}, nil
			},
		},
	}
}

// Return any MessageIDs in the buffer for this UserID
func (s *Implementation) CheckMessages(userID *id.ID, messageID string, ipAddress string) (
	[]string, error) {
	return s.Functions.CheckMessages(userID, messageID, ipAddress)
}

// Returns the message matching the given parameters to the client
func (s *Implementation) GetMessage(userID *id.ID, msgID string, ipAddress string) (
	*pb.Slot, error) {
	return s.Functions.GetMessage(userID, msgID, ipAddress)
}

// Upload a message to the cMix Gateway
func (s *Implementation) PutMessage(message *pb.GatewaySlot, ipAddress string) (*pb.GatewaySlotResponse, error) {
	return s.Functions.PutMessage(message, ipAddress)
}

// Pass-through for Registration Nonce Communication
func (s *Implementation) RequestNonce(message *pb.NonceRequest, ipAddress string) (
	*pb.Nonce, error) {
	return s.Functions.RequestNonce(message, ipAddress)
}

// Pass-through for Registration Nonce Confirmation
func (s *Implementation) ConfirmNonce(message *pb.RequestRegistrationConfirmation,
	ipAddress string) (*pb.RegistrationConfirmation, error) {
	return s.Functions.ConfirmNonce(message, ipAddress)
}

// Ping gateway to ask for users to notify
func (s *Implementation) PollForNotifications(auth *connect.Auth) ([]*id.ID, error) {
	return s.Functions.PollForNotifications(auth)
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

// Client -> Gateway bloom request
func (s *Implementation) RequestBloom(msg *pb.GetBloom) (*pb.GetBloomResponse, error) {
	return s.Functions.RequestBloom(msg)
}
