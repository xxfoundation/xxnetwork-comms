////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/privategrity/comms/mixmessages"
	"runtime/debug"
)

// Handler implementation for the Gateway
type Handler interface {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages(userID uint64, messageID string) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage(userID uint64, msgID string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage(message *pb.CmixMessage) bool
	// ReceiveBatch receives message from a cMix node
	ReceiveBatch(messages *pb.OutputMessages)
}

type implementationFunctions struct {
	// Return any MessageIDs in the buffer for this UserID
	CheckMessages func(userID uint64, messageID string) ([]string, bool)
	// Returns the message matching the given parameters to the client
	GetMessage func(userID uint64, msgID string) (*pb.CmixMessage, bool)
	// Upload a message to the cMix Gateway
	PutMessage func(message *pb.CmixMessage) bool
	// ReceiveBatch receives message from a cMix node
	ReceiveBatch func(messages *pb.OutputMessages)
}

// Implementation allows users of the client library to set the
// functions that implement the node functions
type Implementation struct {
	Functions implementationFunctions
}

func NewImplementation() Handler {
	um := "UNIMPLEMENTED FUNCTION!"
	warn := func(msg string) {
		jww.WARN.Printf(msg)
		jww.WARN.Printf("%v", debug.Stack())
	}
	return Handler(&Implementation{
		Functions: implementationFunctions{
			CheckMessages: func(userID uint64, messageID string) ([]string, bool) {
				warn(um)
				return nil, false
			},
			GetMessage: func(userID uint64, msgID string) (*pb.CmixMessage, bool) {
				warn(um)
				return nil, false
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
func (s *Implementation) CheckMessages(userID uint64, messageID string) (
	[]string, bool) {
	return s.Functions.CheckMessages(userID, messageID)
}

// Returns the message matching the given parameters to the client
func (s *Implementation) GetMessage(userID uint64, msgID string) (
	*pb.CmixMessage, bool) {
	return s.Functions.GetMessage(userID, msgID)
}

// Upload a message to the cMix Gateway
func (s *Implementation) PutMessage(message *pb.CmixMessage) bool {
	return s.Functions.PutMessage(message)
}

// ReceiveBatch receives message from a cMix node
func (s *Implementation) ReceiveBatch(messages *pb.OutputMessages) {
	s.Functions.ReceiveBatch(messages)
}
