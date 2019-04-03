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
	ServerMetrics(*mixmessages.ServerMetricsMessage)

	// Server Interface for starting New Rounds
	NewRound(RoundID string)
	// Server interface for Starting a new round
	StartRound(message *mixmessages.InputMessages)
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo() (int, error)

	// Server Interface for all Internode Comms
	Phase(message *mixmessages.CmixMessage)

	// Server interface for RequestNonceMessage
	RequestNonce(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmNonce(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)
}

type implementationFunctions struct {
	// Server Interface for roundtrip ping
	RoundtripPing func(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	ServerMetrics func(*mixmessages.ServerMetricsMessage)

	// Server Interface for starting New Rounds
	NewRound func(RoundID string)
	// Server interface for Starting a new round
	StartRound func(message *mixmessages.InputMessages)
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func() (int, error)

	// Server Interface for the Internode Messages
	Phase func(message *mixmessages.CmixMessage)

	// Server interface for RequestNonceMessage
	RequestNonce func(salt, Y, P, Q, G,
		hash, R, S []byte) ([]byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmNonce func(hash, R, S []byte) ([]byte,
		[]byte, []byte, []byte, []byte, []byte, []byte, error)
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
			RoundtripPing: func(pingMsg *mixmessages.TimePing) { warn(um) },
			ServerMetrics: func(metMsg *mixmessages.ServerMetricsMessage) { warn(um) },
			NewRound:      func(RoundID string) { warn(um) },
			Phase:         func(m *mixmessages.CmixMessage) { warn(um) },
			StartRound:    func(message *mixmessages.InputMessages) { warn(um) },
			GetRoundBufferInfo: func() (int, error) {
				warn(um)
				return 0, nil
			},

			RequestNonce: func(salt, Y, P, Q, G,
				hash, R, S []byte) ([]byte, error) {
				warn(um)
				return nil, nil
			},
			ConfirmNonce: func(hash, R, S []byte) ([]byte,
				[]byte, []byte, []byte, []byte, []byte, []byte, error) {
				warn(um)
				return nil, nil, nil, nil, nil, nil, nil, nil
			},
		},
	}
}

// Server Interface for roundtrip ping
func (s *Implementation) RoundtripPing(pingMsg *mixmessages.TimePing) {
	s.Functions.RoundtripPing(pingMsg)
}

// Server Interface for ServerMetrics Messages
func (s *Implementation) ServerMetrics(
	metricsMsg *mixmessages.ServerMetricsMessage) {
	s.Functions.ServerMetrics(metricsMsg)
}

// Server Interface for starting New Rounds
func (s *Implementation) NewRound(RoundID string) {
	s.Functions.NewRound(RoundID)
}

// Server Interface for the phase messages
func (s *Implementation) Phase(m *mixmessages.CmixMessage) {
	s.Functions.Phase(m)
}

// Server interface for Starting a new round
func (s *Implementation) StartRound(message *mixmessages.InputMessages) {
	s.Functions.StartRound(message)
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
func (s *Implementation) ConfirmNonce(hash, R, S []byte) ([]byte,
	[]byte, []byte, []byte, []byte, []byte, []byte, error) {
	return s.Functions.ConfirmNonce(hash, R, S)
}
