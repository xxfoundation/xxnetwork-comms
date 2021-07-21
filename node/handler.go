///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains callback interface for server functionality

package node

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/interconnect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/reflection"
	"runtime/debug"
	"strconv"
)

// Server object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
	handler Handler
}

// Starts a new server on the address:port specified by listeningAddr
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNode(id *id.ID, localServer string, interconnectPort int, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	pc, lis, err := connect.StartCommServer(id, localServer,
		certPEMblock, keyPEMblock)
	if err != nil {
		jww.FATAL.Panicf("Unable to start comms server: %+v", err)
	}

	mixmessageServer := Comms{
		ProtoComms: pc,
		handler:    handler,
	}

	// Start up interconnect service
	if interconnectPort != 0 {
		go func() {
			interconnect.StartCMixInterconnect(id, strconv.Itoa(interconnectPort), handler, certPEMblock, keyPEMblock)
		}()
	} else {
		jww.WARN.Printf("Port for consensus not set, interconnect not started")
	}

	go func() {

		// Register GRPC services to the listening address
		mixmessages.RegisterNodeServer(mixmessageServer.LocalServer, &mixmessageServer)
		messages.RegisterGenericServer(mixmessageServer.LocalServer, &mixmessageServer)

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
	UploadUnmixedBatch(server mixmessages.Node_UploadUnmixedBatchServer, auth *connect.Auth) error
	// Server interface for broadcasting when realtime is complete
	FinishRealtime(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo(auth *connect.Auth) (int, error)

	GetMeasure(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error)

	// Server Interface for all Internode Comms
	PostPhase(message *mixmessages.Batch, auth *connect.Auth) error

	StreamPostPhase(server mixmessages.Node_StreamPostPhaseServer, auth *connect.Auth) error

	// Server interface for RequestNonceMessage
	RequestNonce(nonceRequest *mixmessages.NonceRequest, auth *connect.Auth) (*mixmessages.Nonce, error)

	// Server interface for ConfirmNonceMessage
	ConfirmRegistration(requestConfirmation *mixmessages.RequestRegistrationConfirmation,
		auth *connect.Auth) (*mixmessages.RegistrationConfirmation, error)

	// PostPrecompResult interface to finalize both payloads' precomps
	PostPrecompResult(roundID uint64, slots []*mixmessages.Slot, auth *connect.Auth) error

	Poll(msg *mixmessages.ServerPoll, auth *connect.Auth) (*mixmessages.ServerPollResponse, error)

	SendRoundTripPing(ping *mixmessages.RoundTripPing, auth *connect.Auth) error

	AskOnline() error

	RoundError(error *mixmessages.RoundError, auth *connect.Auth) error
	// Consensus node -> cMix node NDF request
	// NOTE: For now cMix nodes serve the NDF to the
	//  consensus nodes, but this will be reversed
	//  once consensus generates the NDF
	GetNDF() (*interconnect.NDF, error)

	// GetPermissioningAddress gets gateway the permissioning server's address
	// from server.
	GetPermissioningAddress() (string, error)

	// Server -> Server initiating multi-party round DH key generation
	StartSharePhase(ri *mixmessages.RoundInfo, auth *connect.Auth) error

	// Server -> Server passing state of multi-party round DH key generation
	SharePhaseRound(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error

	// Server -> Server sending multi-party round DH key
	ShareFinalKey(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error
}

type implementationFunctions struct {
	// Server Interface for starting New Rounds
	CreateNewRound func(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// Server interface for sending a new batch
	UploadUnmixedBatch func(stream mixmessages.Node_UploadUnmixedBatchServer, auth *connect.Auth) error

	// Server interface for finishing the realtime phase
	FinishRealtime func(message *mixmessages.RoundInfo, auth *connect.Auth) error
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func(auth *connect.Auth) (int, error)

	GetMeasure func(message *mixmessages.RoundInfo, auth *connect.Auth) (*mixmessages.RoundMetrics, error)

	// Server Interface for the Internode Messages
	PostPhase func(message *mixmessages.Batch, auth *connect.Auth) error

	// Server interface for internode streaming messages
	StreamPostPhase func(message mixmessages.Node_StreamPostPhaseServer, auth *connect.Auth) error

	// Server interface for RequestNonceMessage
	RequestNonce func(nonceRequest *mixmessages.NonceRequest, auth *connect.Auth) (*mixmessages.Nonce, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration func(requestConfirmation *mixmessages.RequestRegistrationConfirmation,
		auth *connect.Auth) (*mixmessages.RegistrationConfirmation, error)

	// PostPrecompResult interface to finalize both payloads' precomputations
	PostPrecompResult func(roundID uint64,
		slots []*mixmessages.Slot, auth *connect.Auth) error

	Poll func(msg *mixmessages.ServerPoll, auth *connect.Auth) (*mixmessages.ServerPollResponse, error)

	SendRoundTripPing func(ping *mixmessages.RoundTripPing, auth *connect.Auth) error

	AskOnline func() error

	RoundError func(error *mixmessages.RoundError, auth *connect.Auth) error
	// Consensus node -> cMix node NDF request
	// NOTE: For now cMix nodes serve the NDF to the
	//  consensus nodes, but this will be reversed
	//  once consensus generates the NDF
	GetNdf func() (*interconnect.NDF, error)

	// GetPermissioningAddress gets gateway the permissioning server's address
	// from server.
	GetPermissioningAddress func() (string, error)

	// Server -> Server initiating multi-party round DH key generation
	StartSharePhase func(ri *mixmessages.RoundInfo, auth *connect.Auth) error

	// Server -> Server passing state of multi-party round DH key generation
	SharePhaseRound func(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error

	// Server -> Server sending multi-party round DH key
	ShareFinalKey func(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error
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
			PostPhase: func(m *mixmessages.Batch, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			StreamPostPhase: func(message mixmessages.Node_StreamPostPhaseServer, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			UploadUnmixedBatch: func(stream mixmessages.Node_UploadUnmixedBatchServer, auth *connect.Auth) error {
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

			RequestNonce: func(nonceRequest *mixmessages.NonceRequest, auth *connect.Auth) (*mixmessages.Nonce, error) {
				warn(um)
				return &mixmessages.Nonce{}, nil
			},
			ConfirmRegistration: func(requestConfirmation *mixmessages.RequestRegistrationConfirmation,
				auth *connect.Auth) (*mixmessages.RegistrationConfirmation, error) {
				warn(um)
				return &mixmessages.RegistrationConfirmation{}, nil
			},
			PostPrecompResult: func(roundID uint64,
				slots []*mixmessages.Slot, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			Poll: func(msg *mixmessages.ServerPoll, auth *connect.Auth) (*mixmessages.ServerPollResponse, error) {
				warn(um)
				return &mixmessages.ServerPollResponse{}, nil
			},
			SendRoundTripPing: func(ping *mixmessages.RoundTripPing, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			AskOnline: func() error {
				warn(um)
				return nil
			},
			RoundError: func(error *mixmessages.RoundError, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			GetNdf: func() (bytes *interconnect.NDF, err error) {
				warn(um)
				return nil, nil
			},
			GetPermissioningAddress: func() (string, error) {
				warn(um)
				return "", nil
			},
			StartSharePhase: func(roundInfo *mixmessages.RoundInfo, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			SharePhaseRound: func(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error {
				warn(um)
				return nil
			},
			ShareFinalKey: func(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error {
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

func (s *Implementation) UploadUnmixedBatch(stream mixmessages.Node_UploadUnmixedBatchServer,
	auth *connect.Auth) error {
	return s.Functions.UploadUnmixedBatch(stream, auth)
}

// Server Interface for the phase messages
func (s *Implementation) PostPhase(m *mixmessages.Batch, auth *connect.Auth) error {
	return s.Functions.PostPhase(m, auth)
}

// Server Interface for streaming phase messages
func (s *Implementation) StreamPostPhase(m mixmessages.Node_StreamPostPhaseServer, auth *connect.Auth) error {
	return s.Functions.StreamPostPhase(m, auth)
}

// GetRoundBufferInfo returns # of completed precomputations
func (s *Implementation) GetRoundBufferInfo(auth *connect.Auth) (int, error) {
	return s.Functions.GetRoundBufferInfo(auth)
}

// Server interface for RequestNonceMessage
func (s *Implementation) RequestNonce(nonceRequest *mixmessages.NonceRequest, auth *connect.Auth) (*mixmessages.Nonce, error) {
	return s.Functions.RequestNonce(nonceRequest, auth)
}

// Server interface for ConfirmNonceMessage
func (s *Implementation) ConfirmRegistration(requestConfirmation *mixmessages.RequestRegistrationConfirmation,
	auth *connect.Auth) (*mixmessages.RegistrationConfirmation, error) {
	return s.Functions.ConfirmRegistration(requestConfirmation, auth)
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

func (s *Implementation) Poll(msg *mixmessages.ServerPoll, auth *connect.Auth) (*mixmessages.ServerPollResponse, error) {
	return s.Functions.Poll(msg, auth)
}

func (s *Implementation) SendRoundTripPing(ping *mixmessages.RoundTripPing, auth *connect.Auth) error {
	return s.Functions.SendRoundTripPing(ping, auth)
}

// AskOnline blocks until the server is online, or returns an error
func (s *Implementation) AskOnline() error {
	return s.Functions.AskOnline()
}

func (s *Implementation) RoundError(err *mixmessages.RoundError, auth *connect.Auth) error {
	return s.Functions.RoundError(err, auth)
}

// Consensus node -> cMix node NDF request
// NOTE: For now cMix nodes serve the NDF to the
//  consensus nodes, but this will be reversed
//  once consensus generates the NDF
func (s *Implementation) GetNDF() (*interconnect.NDF, error) {
	return s.Functions.GetNdf()
}

// GetPermissioningAddress gets gateway the permissioning server's address from
// server.
func (s *Implementation) GetPermissioningAddress() (string, error) {
	return s.Functions.GetPermissioningAddress()
}

// Server -> Server initiating multi-party round DH key generation
func (s *Implementation) StartSharePhase(ri *mixmessages.RoundInfo, auth *connect.Auth) error {
	return s.Functions.StartSharePhase(ri, auth)
}

// Server -> Server passing state of multi-party round DH key generation
func (s *Implementation) SharePhaseRound(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error {
	return s.Functions.SharePhaseRound(sharedPiece, auth)
}

// Server -> Server sending multi-party round DH final key
func (s *Implementation) ShareFinalKey(sharedPiece *mixmessages.SharePiece, auth *connect.Auth) error {
	return s.Functions.ShareFinalKey(sharedPiece, auth)
}
