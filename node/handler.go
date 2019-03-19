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
	// Server Interface for SetPublicKey
	SetPublicKey(RoundID string, PublicKey []byte)

	// Server Interface for the PrecompDecrypt Messages
	PrecompDecrypt(*mixmessages.PrecompDecryptMessage)
	// Server Interface for the PrecompEncrypt Messages
	PrecompEncrypt(*mixmessages.PrecompEncryptMessage)
	// Server Interface for the PrecompReveal Messages
	PrecompReveal(*mixmessages.PrecompRevealMessage)
	// Server Interface for the PrecompPermute Messages
	PrecompPermute(*mixmessages.PrecompPermuteMessage)
	// Server Interface for the PrecompShare Messages
	PrecompShare(*mixmessages.PrecompShareMessage)
	// Server Interface for the PrecompShareInit Messages
	PrecompShareInit(*mixmessages.PrecompShareInitMessage)
	// Server Interface for the PrecompShareInit Messages
	PrecompShareCompare(*mixmessages.PrecompShareCompareMessage)
	// Server Interface for the PrecompShareConfirm Messages
	PrecompShareConfirm(*mixmessages.PrecompShareConfirmMessage)
	// Server Interface for the RealtimeDecrypt Messages
	RealtimeDecrypt(*mixmessages.RealtimeDecryptMessage)
	// Server Interface for the RealtimeEncrypt Messages
	RealtimeEncrypt(*mixmessages.RealtimeEncryptMessage)
	// Server Interface for the RealtimePermute Messages
	RealtimePermute(*mixmessages.RealtimePermuteMessage)

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
	// Server Interface for SetPublicKey
	SetPublicKey func(RoundID string, PublicKey []byte)

	// Server Interface for the PrecompDecrypt Messages
	PrecompDecrypt func(*mixmessages.PrecompDecryptMessage)
	// Server Interface for the PrecompEncrypt Messages
	PrecompEncrypt func(*mixmessages.PrecompEncryptMessage)
	// Server Interface for the PrecompReveal Messages
	PrecompReveal func(*mixmessages.PrecompRevealMessage)
	// Server Interface for the PrecompPermute Messages
	PrecompPermute func(*mixmessages.PrecompPermuteMessage)
	// Server Interface for the PrecompShare Messages
	PrecompShare func(*mixmessages.PrecompShareMessage)
	// Server Interface for the PrecompShareInit Messages
	PrecompShareInit func(*mixmessages.PrecompShareInitMessage)
	// Server Interface for the PrecompShareInit Messages
	PrecompShareCompare func(*mixmessages.PrecompShareCompareMessage)
	// Server Interface for the PrecompShareConfirm Messages
	PrecompShareConfirm func(*mixmessages.PrecompShareConfirmMessage)
	// Server Interface for the RealtimeDecrypt Messages
	RealtimeDecrypt func(*mixmessages.RealtimeDecryptMessage)
	// Server Interface for the RealtimeEncrypt Messages
	RealtimeEncrypt func(*mixmessages.RealtimeEncryptMessage)
	// Server Interface for the RealtimePermute Messages
	RealtimePermute func(*mixmessages.RealtimePermuteMessage)

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
			RoundtripPing:    func(pingMsg *mixmessages.TimePing) { warn(um) },
			ServerMetrics:    func(metMsg *mixmessages.ServerMetricsMessage) { warn(um) },
			NewRound:         func(RoundID string) { warn(um) },
			SetPublicKey:     func(RoundID string, PublicKey []byte) { warn(um) },
			PrecompDecrypt:   func(m *mixmessages.PrecompDecryptMessage) { warn(um) },
			PrecompEncrypt:   func(m *mixmessages.PrecompEncryptMessage) { warn(um) },
			PrecompReveal:    func(m *mixmessages.PrecompRevealMessage) { warn(um) },
			PrecompPermute:   func(m *mixmessages.PrecompPermuteMessage) { warn(um) },
			PrecompShare:     func(m *mixmessages.PrecompShareMessage) { warn(um) },
			PrecompShareInit: func(m *mixmessages.PrecompShareInitMessage) { warn(um) },
			PrecompShareCompare: func(m *mixmessages.PrecompShareCompareMessage) {
				warn(um)
			},
			PrecompShareConfirm: func(m *mixmessages.PrecompShareConfirmMessage) {
				warn(um)
			},

			RealtimeDecrypt: func(m *mixmessages.RealtimeDecryptMessage) { warn(um) },
			RealtimeEncrypt: func(m *mixmessages.RealtimeEncryptMessage) { warn(um) },
			RealtimePermute: func(m *mixmessages.RealtimePermuteMessage) { warn(um) },
			StartRound:      func(message *mixmessages.InputMessages) { warn(um) },
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

// Server Interface for SetPublicKey
func (s *Implementation) SetPublicKey(RoundID string, PublicKey []byte) {
	s.Functions.SetPublicKey(RoundID, PublicKey)
}

// Server Interface for the PrecompDecrypt Messages
func (s *Implementation) PrecompDecrypt(m *mixmessages.PrecompDecryptMessage) {
	s.Functions.PrecompDecrypt(m)
}

// Server Interface for the PrecompEncrypt Messages
func (s *Implementation) PrecompEncrypt(m *mixmessages.PrecompEncryptMessage) {
	s.Functions.PrecompEncrypt(m)
}

// Server Interface for the PrecompReveal Messages
func (s *Implementation) PrecompReveal(m *mixmessages.PrecompRevealMessage) {
	s.Functions.PrecompReveal(m)
}

// Server Interface for the PrecompPermute Messages
func (s *Implementation) PrecompPermute(m *mixmessages.PrecompPermuteMessage) {
	s.Functions.PrecompPermute(m)
}

// Server Interface for the PrecompShare Messages
func (s *Implementation) PrecompShare(m *mixmessages.PrecompShareMessage) {
	s.Functions.PrecompShare(m)
}

// Server Interface for the PrecompShareInit Messages
func (s *Implementation) PrecompShareInit(
	m *mixmessages.PrecompShareInitMessage) {
	s.Functions.PrecompShareInit(m)
}

// Server Interface for the PrecompShareInit Messages
func (s *Implementation) PrecompShareCompare(
	m *mixmessages.PrecompShareCompareMessage) {
	s.Functions.PrecompShareCompare(m)
}

// Server Interface for the PrecompShareConfirm Messages
func (s *Implementation) PrecompShareConfirm(
	m *mixmessages.PrecompShareConfirmMessage) {
	s.Functions.PrecompShareConfirm(m)
}

// Server Interface for the RealtimeDecrypt Messages
func (s *Implementation) RealtimeDecrypt(
	m *mixmessages.RealtimeDecryptMessage) {
	s.Functions.RealtimeDecrypt(m)
}

// Server Interface for the RealtimeEncrypt Messages
func (s *Implementation) RealtimeEncrypt(
	m *mixmessages.RealtimeEncryptMessage) {
	s.Functions.RealtimeEncrypt(m)
}

// Server Interface for the RealtimePermute Messages
func (s *Implementation) RealtimePermute(
	m *mixmessages.RealtimePermuteMessage) {
	s.Functions.RealtimePermute(m)
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
