////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for server functionality

package node

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
)

// Server object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

// Starts a new server on the address:port specified by listeningAddr
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNode(id, localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Printf("Unable to start comms server: %+v", err)
	}

	mixmessageServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	go func() {
		// Register GRPC services to the listening address
		mixmessages.RegisterNodeServer(mixmessageServer.LocalServer, &mixmessageServer)
		mixmessages.RegisterGenericServer(mixmessageServer.LocalServer, &mixmessageServer)

		// Register reflection service on gRPC server.
		reflection.Register(mixmessageServer.LocalServer)
		if err := mixmessageServer.LocalServer.Serve(lis); err != nil {
			jww.FATAL.Panicf("Failed to serve: %+v",
				errors.New(err.Error()))
		}
		jww.INFO.Printf("Shutting down node server listener: %s", lis)
	}()

	return &mixmessageServer
}

type Handler interface {
	// Server interface for starting New Rounds
	CreateNewRound(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// Server interface for sending a new batch
	PostNewBatch(message *mixmessages.Batch, auth *connect.Auth) error
	// Server interface for broadcasting when realtime is complete
	FinishRealtime(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo(auth *connect.Auth) (int, error)

	GetMeasure(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error)

	// Server Interface for all Internode Comms
	PostPhase(message *mixmessages.Batch, auth *connect.Auth)

	StreamPostPhase(server mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey(message *mixmessages.RoundPublicKey, auth *connect.Auth)

	// Server interface for RequestNonceMessage
	RequestNonce(salt []byte, RSAPubKey string, DHPubKey,
		RSASignedByRegistration, DHSignedByClientRSA []byte, auth *connect.Auth) ([]byte, []byte, error)

	// Server interface for ConfirmNonceMessage
	ConfirmRegistration(UserID []byte, Signature []byte, auth *connect.Auth) ([]byte, error)

	// PostPrecompResult interface to finalize both payloads' precomps
	PostPrecompResult(roundID uint64, slots []*mixmessages.Slot, auth *connect.Auth) error

	// GetCompletedBatch: gateway uses completed batch from the server
	GetCompletedBatch(auth *connect.Auth) (*mixmessages.Batch, error)

	PollNdf(ping *mixmessages.Ping, auth *connect.Auth) (*mixmessages.GatewayNdf, error)

	SendRoundTripPing(ping *mixmessages.RoundTripPing, auth *connect.Auth) error

	AskOnline(ping *mixmessages.Ping, auth *connect.Auth) error
}

type implementationFunctions struct {
	// Server Interface for starting New Rounds
	CreateNewRound func(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// Server interface for sending a new batch
	PostNewBatch func(message *mixmessages.Batch, auth *connect.Auth) error
	// Server interface for finishing the realtime phase
	FinishRealtime func(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func(auth *connect.Auth) (int, error)

	GetMeasure func(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error)

	// Server Interface for the Internode Messages
	PostPhase func(message *mixmessages.Batch, auth *connect.Auth)

	// Server interface for internode streaming messages
	StreamPostPhase func(message mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey func(message *mixmessages.RoundPublicKey, auth *connect.Auth)

	// Server interface for RequestNonceMessage
	RequestNonce func(salt []byte, RSAPubKey string, DHPubKey,
		RSASigFromReg, RSASigDH []byte, auth *connect.Auth) ([]byte, []byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration func(UserID, Signature []byte, auth *connect.Auth) ([]byte, error)

	// PostPrecompResult interface to finalize both payloads' precomputations
	PostPrecompResult func(roundID uint64,
		slots []*mixmessages.Slot, auth *connect.Auth) error

	GetCompletedBatch func(auth *connect.Auth) (*mixmessages.Batch, error)

	PollNdf func(ping *mixmessages.Ping, auth *connect.Auth) (*mixmessages.GatewayNdf, error)

	SendRoundTripPing func(ping *mixmessages.RoundTripPing, auth *connect.Auth) error

	AskOnline func(ping *mixmessages.Ping, auth *connect.Auth) error
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// Below is the Implementation implementation, which calls the
// function matching the variable in the structure.

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
func NewImplementation() *Implementation {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%s", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			CreateNewRound: func(m *mixmessages.RoundInfo, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			PostPhase: func(m *mixmessages.Batch, auth *connect.Auth) {
				warn(um)
			},
			StreamPostPhase: func(message mixmessages.Node_StreamPostPhaseServer) error {
				warn(um)
				return nil
			},
			PostRoundPublicKey: func(message *mixmessages.RoundPublicKey, auth *connect.Auth) {
				warn(um)
			},
			PostNewBatch: func(message *mixmessages.Batch, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			FinishRealtime: func(message *mixmessages.RoundInfo, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			GetMeasure: func(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error) {
				warn(um)
				return nil, nil
			},
			GetRoundBufferInfo: func(auth *connect.Auth) (int, error) {
				warn(um)
				return 0, nil
			},

			RequestNonce: func(salt []byte, RSAPubKey string, DHPubKey,
				RSASig, RSASigDH []byte, auth *connect.Auth) ([]byte, []byte, error) {
				warn(um)
				return nil, nil, nil
			},
			ConfirmRegistration: func(UserID, Signature []byte, auth *connect.Auth) ([]byte, error) {
				warn(um)
				return nil, nil
			},
			PostPrecompResult: func(roundID uint64,
				slots []*mixmessages.Slot, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			GetCompletedBatch: func(auth *connect.Auth) (batch *mixmessages.Batch, e error) {
				warn(um)
				return &mixmessages.Batch{}, nil
			},
			PollNdf: func(ping *mixmessages.Ping, auth *connect.Auth) (certs *mixmessages.GatewayNdf,
				e error) {
				warn(um)
				return &mixmessages.GatewayNdf{}, nil
			},
			SendRoundTripPing: func(ping *mixmessages.RoundTripPing, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			AskOnline: func(ping *mixmessages.Ping, auth *connect.Auth) error {
				warn(um)
				return nil
			},
		},
	}
}

// Server Interface for starting New Rounds
func (s *Implementation) CreateNewRound(msg *mixmessages.RoundInfo, auth *connect.Auth) error {
	return s.Functions.CreateNewRound(msg, auth)
}

func (s *Implementation) PostNewBatch(msg *mixmessages.Batch, auth *connect.Auth) error {
	return s.Functions.PostNewBatch(msg, auth)
}

// Server Interface for the phase messages
func (s *Implementation) PostPhase(m *mixmessages.Batch, auth *connect.Auth) {
	s.Functions.PostPhase(m, auth)
}

// Server Interface for streaming phase messages
func (s *Implementation) StreamPostPhase(m mixmessages.Node_StreamPostPhaseServer) error {
	return s.Functions.StreamPostPhase(m)
}

// Server Interface for the share message
func (s *Implementation) PostRoundPublicKey(message *mixmessages.
	RoundPublicKey, auth *connect.Auth) {
	s.Functions.PostRoundPublicKey(message, auth)
}

// GetRoundBufferInfo returns # of completed precomputations
func (s *Implementation) GetRoundBufferInfo(auth *connect.Auth) (int, error) {
	return s.Functions.GetRoundBufferInfo(auth)
}

// Server interface for RequestNonceMessage
func (s *Implementation) RequestNonce(salt []byte, RSAPubKey string, DHPubKey,
	RSASigFromReg, RSASigDH []byte, auth *connect.Auth) ([]byte, []byte, error) {
	return s.Functions.RequestNonce(salt, RSAPubKey, DHPubKey, RSASigFromReg, RSASigDH, auth)
}

// Server interface for ConfirmNonceMessage
func (s *Implementation) ConfirmRegistration(UserID, Signature []byte, auth *connect.Auth) ([]byte, error) {
	return s.Functions.ConfirmRegistration(UserID, Signature, auth)
}

// PostPrecompResult interface to finalize both payloads' precomputations
func (s *Implementation) PostPrecompResult(roundID uint64,
	slots []*mixmessages.Slot, auth *connect.Auth) error {
	return s.Functions.PostPrecompResult(roundID, slots, auth)
}

func (s *Implementation) FinishRealtime(message *mixmessages.RoundInfo, auth *connect.Auth) error {
	return s.Functions.FinishRealtime(message, auth)
}

func (s *Implementation) GetMeasure(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error) {
	return s.Functions.GetMeasure(message, auth)
}

// Implementation of the interface using the function in the struct
func (s *Implementation) GetCompletedBatch(auth *connect.Auth) (*mixmessages.Batch, error) {
	return s.Functions.GetCompletedBatch(auth)
}

func (s *Implementation) PollNdf(ping *mixmessages.Ping, auth *connect.Auth) (*mixmessages.
	GatewayNdf, error) {
	return s.Functions.PollNdf(ping, auth)
}

func (s *Implementation) SendRoundTripPing(ping *mixmessages.RoundTripPing, auth *connect.Auth) error {
	return s.Functions.SendRoundTripPing(ping, auth)
}

// AskOnline blocks until the server is online, or returns an error
func (s *Implementation) AskOnline(ping *mixmessages.Ping, auth *connect.Auth) error {
	return s.Functions.AskOnline(ping, auth)
}
