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
	// Server interface for starting New Rounds
	CreateNewRound(message *mixmessages.RoundInfo) error
	// Server interface for sending a new batch
	PostNewBatch(message *mixmessages.Batch) error
	// Server interface for broadcasting when realtime is complete
	FinishRealtime(message *mixmessages.RoundInfo) error
	// GetRoundBufferInfo returns # of available precomputations
	GetRoundBufferInfo() (int, error)

	GetMeasure(message *mixmessages.RoundInfo) (*mixmessages.RoundMetrics, error)

	// Server Interface for all Internode Comms
	PostPhase(message *mixmessages.Batch)

	StreamPostPhase(server mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey(message *mixmessages.RoundPublicKey)

	// Server interface for RequestNonceMessage
	RequestNonce(salt []byte, RSAPubKey string, DHPubKey,
		RSASignedByRegistration, DHSignedByClientRSA []byte) ([]byte, []byte, error)

	// Server interface for ConfirmNonceMessage
	ConfirmRegistration(UserID []byte, Signature []byte) ([]byte, error)

	// PostPrecompResult interface to finalize both payloads' precomps
	PostPrecompResult(roundID uint64, slots []*mixmessages.Slot) error

	// GetCompletedBatch: gateway uses completed batch from the server
	GetCompletedBatch() (*mixmessages.Batch, error)

	// DownloadTopology: Obtains network topology from permissioning server
	DownloadTopology(info *MessageInfo, topology *mixmessages.NodeTopology)

	GetSignedCert(ping *mixmessages.Ping) (*mixmessages.SignedCerts, error)
}

type implementationFunctions struct {
	// Server Interface for starting New Rounds
	CreateNewRound func(message *mixmessages.RoundInfo) error
	// Server interface for sending a new batch
	PostNewBatch func(message *mixmessages.Batch) error
	// Server interface for finishing the realtime phase
	FinishRealtime func(message *mixmessages.RoundInfo) error
	// GetRoundBufferInfo returns # of available precomputations completed
	GetRoundBufferInfo func() (int, error)

	GetMeasure func(message *mixmessages.RoundInfo) (*mixmessages.RoundMetrics, error)

	// Server Interface for the Internode Messages
	PostPhase func(message *mixmessages.Batch)

	// Server interface for internode streaming messages
	StreamPostPhase func(message mixmessages.Node_StreamPostPhaseServer) error

	// Server interface for share broadcast
	PostRoundPublicKey func(message *mixmessages.RoundPublicKey)

	// Server interface for RequestNonceMessage
	RequestNonce func(salt []byte, RSAPubKey string, DHPubKey,
		RSASigFromReg, RSASigDH []byte) ([]byte, []byte, error)
	// Server interface for ConfirmNonceMessage
	ConfirmRegistration func(UserID, Signature []byte) ([]byte, error)

	// PostPrecompResult interface to finalize both payloads' precomputations
	PostPrecompResult func(roundID uint64,
		slots []*mixmessages.Slot) error

	GetCompletedBatch func() (*mixmessages.Batch, error)

	// DownloadTopology: Obtains network topology from permissioning server
	DownloadTopology func(info *MessageInfo, topology *mixmessages.NodeTopology)

	GetSignedCert func(ping *mixmessages.Ping) (*mixmessages.SignedCerts, error)
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
			GetMeasure: func(message *mixmessages.RoundInfo) (*mixmessages.RoundMetrics, error) {
				warn(um)
				return nil, nil
			},
			GetRoundBufferInfo: func() (int, error) {
				warn(um)
				return 0, nil
			},

			RequestNonce: func(salt []byte, RSAPubKey string, DHPubKey,
				RSASig, RSASigDH []byte) ([]byte, []byte, error) {
				warn(um)
				return nil, nil, nil
			},
			ConfirmRegistration: func(UserID, Signature []byte) ([]byte, error) {
				warn(um)
				return nil, nil
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
			DownloadTopology: func(info *MessageInfo, topology *mixmessages.NodeTopology) {
				warn(um)
			},
			GetSignedCert: func(ping *mixmessages.Ping) (certs *mixmessages.SignedCerts, e error) {
				warn(um)
				return &mixmessages.SignedCerts{}, nil
			},
		},
	}
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
func (s *Implementation) RequestNonce(salt []byte, RSAPubKey string, DHPubKey,
	RSASigFromReg, RSASigDH []byte) ([]byte, []byte, error) {
	return s.Functions.RequestNonce(salt, RSAPubKey, DHPubKey, RSASigFromReg, RSASigDH)
}

// Server interface for ConfirmNonceMessage
func (s *Implementation) ConfirmRegistration(UserID, Signature []byte) ([]byte, error) {
	return s.Functions.ConfirmRegistration(UserID, Signature)
}

// PostPrecompResult interface to finalize both payloads' precomputations
func (s *Implementation) PostPrecompResult(roundID uint64,
	slots []*mixmessages.Slot) error {
	return s.Functions.PostPrecompResult(roundID, slots)
}

func (s *Implementation) FinishRealtime(message *mixmessages.RoundInfo) error {
	return s.Functions.FinishRealtime(message)
}

func (s *Implementation) GetMeasure(message *mixmessages.RoundInfo) (*mixmessages.RoundMetrics, error) {
	return s.Functions.GetMeasure(message)
}

// Implementation of the interface using the function in the struct
func (s *Implementation) GetCompletedBatch() (*mixmessages.Batch, error) {
	return s.Functions.GetCompletedBatch()
}

// Obtains network topology from permissioning server
func (s *Implementation) DownloadTopology(info *MessageInfo, topology *mixmessages.NodeTopology) {
	s.Functions.DownloadTopology(info, topology)
}

func (s *Implementation) GetSignedCert(ping *mixmessages.Ping) (*mixmessages.SignedCerts, error) {
	return s.Functions.GetSignedCert(ping)
}
