////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for server functionality

package node

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/mixmessages"
	"runtime/debug"
)

type ServerHandler interface {
	// Server interface for round trip ping
	RoundtripPing(*mixmessages.TimePing)
	// Server interface for ServerMetrics Messages
	GetServerMetrics(*mixmessages.ServerMetrics)

	// Server interface for starting New Rounds
	CreateNewRound(message *mixmessages.RoundInfo) error
	// Server interface for sending a new batch
	PostNewBatch(message *mixmessages.Batch) error
	// Server interface for broadcasting when realtime is complete
	FinishRealtime(message *mixmessages.RoundInfo) error
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo() (int, error)

	GetMeasure(message *mixmessages.RoundInfo) error

	// Server Interface for all Internode Comms
	PostPhase(message *mixmessages.Batch)

	StreamPostPhase(server mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey(message *mixmessages.RoundPublicKey)

	// Server interface for RequestNonceMessage
	RequestNonce(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)

	// PostPrecompResult interface to finalize message and AD precomps
	PostPrecompResult(roundID uint64, slots []*mixmessages.Slot) error

	// GetCompletedBatch: gateway uses completed batch from the server
	GetCompletedBatch() (*mixmessages.Batch, error)

	// DownloadTopology: Obtains network topology from permissioning server
	DownloadTopology(topology *mixmessages.NodeTopology)
}

type implementationFunctions struct {
	// Server Interface for roundtrip ping
	RoundtripPing func(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	GetServerMetrics func(*mixmessages.ServerMetrics)

	// Server Interface for starting New Rounds
	CreateNewRound func(message *mixmessages.RoundInfo) error
	// Server interface for sending a new batch
	PostNewBatch func(message *mixmessages.Batch) error
	// Server interface for finishing the realtime phase
	FinishRealtime func(message *mixmessages.RoundInfo) error
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func() (int, error)

	GetMeasure func(message *mixmessages.RoundInfo) error

	// Server Interface for the Internode Messages
	PostPhase func(message *mixmessages.Batch)

	// Server interface for internode streaming messages
	StreamPostPhase func(message mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey func(message *mixmessages.RoundPublicKey)

	// Server interface for RequestNonceMessage
	RequestNonce func(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration func(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)

	// PostPrecompResult interface to finalize message and AD precomps
	PostPrecompResult func(roundID uint64,
		slots []*mixmessages.Slot) error

	GetCompletedBatch func() (*mixmessages.Batch, error)

	// DownloadTopology: Obtains network topology from permissioning server
	DownloadTopology func(topology *mixmessages.NodeTopology)
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
		jww.WARN.Printf("%v", debug.Stack())
	}
	return &Implementation{
		Functions: implementationFunctions{
			RoundtripPing: func(p *mixmessages.TimePing) {
				warn(um)
			},
			GetServerMetrics: func(m *mixmessages.ServerMetrics) {
				warn(um)
			},
			CreateNewRound: func(m *mixmessages.RoundInfo) error {
				warn(um)
				return nil
			},
			PostPhase: func(m *mixmessages.Batch) {
				warn(um)
			},
			StreamPostPhase: func(message mixmessages.Node_StreamPostPhaseServer) error {
				warn(um)
				return nil
			},
			PostRoundPublicKey: func(message *mixmessages.RoundPublicKey) {
				warn(um)
			},
			PostNewBatch: func(message *mixmessages.Batch) error {
				warn(um)
				return nil
			},
			FinishRealtime: func(message *mixmessages.RoundInfo) error {
				warn(um)
				return nil
			},
			GetMeasure: func(message *mixmessages.RoundInfo) error {
				warn(um)
				return nil
			},
			GetRoundBufferInfo: func() (int, error) {
				warn(um)
				return 0, nil
			},

			RequestNonce: func(salt, Y, P, Q, G,
				hash, R, S []byte) ([]byte, error) {
				warn(um)
				return nil, nil
			},
			ConfirmRegistration: func(hash, R, S []byte) (
				[]byte, []byte, []byte, []byte, []byte,
				[]byte, []byte, error) {
				warn(um)
				return nil, nil, nil, nil, nil, nil, nil, nil
			},
			PostPrecompResult: func(roundID uint64,
				slots []*mixmessages.Slot) error {
				warn(um)
				return nil
			},
			GetCompletedBatch: func() (batch *mixmessages.Batch, e error) {
				warn(um)
				return &mixmessages.Batch{}, nil
			},
			DownloadTopology: func(topology *mixmessages.NodeTopology) {
				warn(um)
			},
		},
	}
}

// Server Interface for roundtrip ping
func (s *Implementation) RoundtripPing(pingMsg *mixmessages.TimePing) {
	s.Functions.RoundtripPing(pingMsg)
}

// Server Interface for ServerMetrics Messages
func (s *Implementation) GetServerMetrics(
	metricsMsg *mixmessages.ServerMetrics) {
	s.Functions.GetServerMetrics(metricsMsg)
}

// Server Interface for starting New Rounds
func (s *Implementation) CreateNewRound(msg *mixmessages.RoundInfo) error {
	return s.Functions.CreateNewRound(msg)
}

func (s *Implementation) PostNewBatch(msg *mixmessages.Batch) error {
	return s.Functions.PostNewBatch(msg)
}

// Server Interface for the phase messages
func (s *Implementation) PostPhase(m *mixmessages.Batch) {
	s.Functions.PostPhase(m)
}

// Server Interface for streaming phase messages
func (s *Implementation) StreamPostPhase(m mixmessages.Node_StreamPostPhaseServer) error {
	return s.Functions.StreamPostPhase(m)
}

// Server Interface for the share message
func (s *Implementation) PostRoundPublicKey(message *mixmessages.
	RoundPublicKey) {
	s.Functions.PostRoundPublicKey(message)
}

// GetRoundBufferInfo returns # of completed precomputations
func (s *Implementation) GetRoundBufferInfo() (int, error) {
	return s.Functions.GetRoundBufferInfo()
}

// Server interface for RequestNonceMessage
func (s *Implementation) RequestNonce(salt, Y, P, Q, G,
	hash, R, S []byte) ([]byte, error) {
	return s.Functions.RequestNonce(salt, Y, P, Q, G, hash, R, S)
}

// Server interface for ConfirmNonceMessage
func (s *Implementation) ConfirmRegistration(hash, R, S []byte) ([]byte,
	[]byte, []byte, []byte, []byte, []byte, []byte, error) {
	return s.Functions.ConfirmRegistration(hash, R, S)
}

// PostPrecompResult interface to finalize message and AD precomps
func (s *Implementation) PostPrecompResult(roundID uint64,
	slots []*mixmessages.Slot) error {
	return s.Functions.PostPrecompResult(roundID, slots)
}

func (s *Implementation) FinishRealtime(message *mixmessages.RoundInfo) error {
	return s.Functions.FinishRealtime(message)
}

func (s *Implementation) GetMeasure(message *mixmessages.RoundInfo) error {
	return s.Functions.GetMeasure(message)
}

// Implementation of the interface using the function in the struct
func (s *Implementation) GetCompletedBatch() (*mixmessages.Batch, error) {
	return s.Functions.GetCompletedBatch()
}

// Obtains network topology from permissioning server
func (s *Implementation) DownloadTopology(topology *mixmessages.NodeTopology) {
	s.Functions.DownloadTopology(topology)
}
