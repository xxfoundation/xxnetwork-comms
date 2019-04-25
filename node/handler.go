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
	// Server Interface for round trip ping
	RoundtripPing(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	GetServerMetrics(*mixmessages.ServerMetrics)

	// Server Interface for starting New Rounds
	CreateNewRound(RoundID uint64)
	// Server interface for Starting starting realtime
	StartRealtime(message *mixmessages.Input)
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo() (int, error)

	// Server Interface for all Internode Comms
	RunPhase(message *mixmessages.Batch)

	// Server interface for RequestNonceMessage
	RequestNonce(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)

	// FinishPrecomputation interface to finalize message and AD precomps
	FinishPrecomputation(roundID uint64, slots []*mixmessages.Slot) error
}

type implementationFunctions struct {
	// Server Interface for roundtrip ping
	RoundtripPing func(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	GetServerMetrics func(*mixmessages.ServerMetrics)

	// Server Interface for starting New Rounds
	CreateNewRound func(RoundID uint64)
	// Server interface for Starting the realtime phase
	StartRealtime func(message *mixmessages.Input)
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func() (int, error)

	// Server Interface for the Internode Messages
	RunPhase func(message *mixmessages.Batch)

	// Server interface for RequestNonceMessage
	RequestNonce func(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration func(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)

	// FinishPrecomputation interface to finalize message and AD precomps
	FinishPrecomputation func(roundID uint64,
		slots []*mixmessages.Slot) error
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
			CreateNewRound: func(RoundID uint64) { warn(um) },
			RunPhase: func(m *mixmessages.Batch) {
				warn(um)
			},
			StartRealtime: func(message *mixmessages.Input) {
				warn(um)
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
			FinishPrecomputation: func(roundID uint64,
				slots []*mixmessages.Slot) error {
				warn(um)
				return nil
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
func (s *Implementation) CreateNewRound(RoundID uint64) {
	s.Functions.CreateNewRound(RoundID)
}

// Server Interface for the phase messages
func (s *Implementation) RunPhase(m *mixmessages.Batch) {
	s.Functions.RunPhase(m)
}

// Server interface for Starting a new round
func (s *Implementation) StartRealtime(message *mixmessages.Input) {
	s.Functions.StartRealtime(message)
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

// FinishPrecomputation interface to finalize message and AD precomps
func (s *Implementation) FinishPrecomputation(roundID uint64,
	slots []*mixmessages.Slot) error {
	return s.Functions.FinishPrecomputation(roundID, slots)
}
