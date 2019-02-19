////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains callback interface for gateway functionality

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"runtime/debug"
)

// Handler interface for the Gateway
type Handler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userID *id.User, messageID string) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage(userID *id.User, msgID string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage(message *pb.CmixMessage) bool
	// Receives a batch of messages from a server
	ReceiveBatch(messages *pb.OutputMessages)
}

// Handler implementation for the Gateway
type implementationFunctions struct {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages func(userID *id.User, messageID string) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage func(userID *id.User, msgID string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage func(message *pb.CmixMessage) bool
	// Receives a batch of messages from a server
	ReceiveBatch func(messages *pb.OutputMessages)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

// Creates and returns a new Handler interface
func NewImplementation() Handler {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%v", debug.Stack())
	}
	return Handler(&Implementation{
		Functions: implementationFunctions{
			CheckMessages: func(userID *id.User, messageID string) ([]string,
				bool) {
				warn(um)
				return nil, false
			},
			GetMessage: func(userID *id.User, msgID string) (*pb.CmixMessage,
				bool) {
				warn(um)
				return &pb.CmixMessage{}, false
			},
			PutMessage: func(message *pb.CmixMessage) bool {
				warn(um)
				return false
			},
			ReceiveBatch: func(messages *pb.OutputMessages) { warn(um) },
		},
	})
}

// Return any MessageIDs in the buffer for this UserID
func (s *Implementation) CheckMessages(userID *id.User, messageID string) (
	[]string, bool) {
	return s.Functions.CheckMessages(userID, messageID)
}

// Returns the message matching the given parameters to the client
func (s *Implementation) GetMessage(userID *id.User, msgID string) (
	*pb.CmixMessage, bool) {
	return s.Functions.GetMessage(userID, msgID)
}

// Upload a message to the cMix Gateway
func (s *Implementation) PutMessage(message *pb.CmixMessage) bool {
	return s.Functions.PutMessage(message)
}

// Receives a batch of messages from a server
func (s *Implementation) ReceiveBatch(messages *pb.OutputMessages) {
	s.Functions.ReceiveBatch(messages)
}
