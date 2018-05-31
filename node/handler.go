////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/privategrity/comms/mixmessages"
	"runtime/debug"
)

type ServerHandler interface {
	// Server Interface for roundtrip ping
	RoundtripPing(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	ServerMetrics(*mixmessages.ServerMetricsMessage)

	// Server Interface for starting New Rounds
	NewRound(RoundID string)
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

	// Server Interface for responding to ClientPoll Messages
	ClientPoll(*mixmessages.ClientPollMessage) *mixmessages.CmixMessage
	// Server interface for ReceiveMessageFromClient
	ReceiveMessageFromClient(*mixmessages.CmixMessage)
	// Server interface for responding to contact list requests from the client
	RequestContactList(*mixmessages.ContactPoll) *mixmessages.ContactMessage

	// Server interface for upserting a new user
	UserUpsert(message *mixmessages.UpsertUserMessage)
	// Check the registration status of a specific user
	PollRegistrationStatus(message *mixmessages.RegistrationPoll) *mixmessages.RegistrationConfirmation
	// Server interface for Starting a new round
	StartRound(message *mixmessages.InputMessages)
}

type implementationFunctions struct {
	// Server Interface for roundtrip ping
	RoundtripPing func(*mixmessages.TimePing)
	// Server Interface for ServerMetrics Messages
	ServerMetrics func(*mixmessages.ServerMetricsMessage)

	// Server Interface for starting New Rounds
	NewRound func(RoundID string)
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

	// Server Interface for responding to ClientPoll Messages
	ClientPoll func(*mixmessages.ClientPollMessage) *mixmessages.CmixMessage
	// Server interface for ReceiveMessageFromClient
	ReceiveMessageFromClient func(*mixmessages.CmixMessage)
	// Server interface for responding to contact list requests from the client
	RequestContactList func(*mixmessages.ContactPoll) *mixmessages.ContactMessage

	// Server interface for upserting a new user
	UserUpsert func(message *mixmessages.UpsertUserMessage)
	// Check the registration status of a specific user
	PollRegistrationStatus func(message *mixmessages.RegistrationPoll) *mixmessages.RegistrationConfirmation
	// Server interface for Starting a new round
	StartRound func(message *mixmessages.InputMessages)
}

type Implementation struct {
	Functions implementationFunctions
}

// Below is the Implementation implementation, which calls the
// function matching the variable in the structure.

// NewImplementation returns a Implementation struct with all of the
// function pointers returning nothing and printing an error.
func NewImplementation() ServerHandler {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%v", debug.Stack())
	}
	return ServerHandler(&Implementation{
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
			ClientPoll: func(m *mixmessages.ClientPollMessage) *mixmessages.CmixMessage {
				warn(um)
				return nil
			},
			ReceiveMessageFromClient: func(m *mixmessages.CmixMessage) { warn(um) },
			RequestContactList: func(m *mixmessages.ContactPoll) *mixmessages.ContactMessage {
				warn(um)
				return nil
			},
			UserUpsert: func(message *mixmessages.UpsertUserMessage) { warn(um) },
			PollRegistrationStatus: func(message *mixmessages.RegistrationPoll) *mixmessages.RegistrationConfirmation {
				warn(um)
				return nil
			},
			StartRound: func(message *mixmessages.InputMessages) { warn(um) },
		},
	})
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

// Server Interface for responding to ClientPoll Messages
func (s *Implementation) ClientPoll(m *mixmessages.ClientPollMessage) *mixmessages.CmixMessage {
	return s.Functions.ClientPoll(m)
}

// Server interface for ReceiveMessageFromClient
func (s *Implementation) ReceiveMessageFromClient(m *mixmessages.CmixMessage) {
	s.Functions.ReceiveMessageFromClient(m)
}

// Server interface for responding to contact list requests from the client
func (s *Implementation) RequestContactList(m *mixmessages.ContactPoll) *mixmessages.ContactMessage {
	return s.Functions.RequestContactList(m)
}

// Server interface for upserting a new user
func (s *Implementation) UserUpsert(message *mixmessages.UpsertUserMessage) {
	s.Functions.UserUpsert(message)
}

// Check the registration status of a specific user
func (s *Implementation) PollRegistrationStatus(
	message *mixmessages.RegistrationPoll) *mixmessages.RegistrationConfirmation {
	return s.Functions.PollRegistrationStatus(message)
}

// Server interface for Starting a new round
func (s *Implementation) StartRound(message *mixmessages.InputMessages) {
	s.Functions.StartRound(message)
}
