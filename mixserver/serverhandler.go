// serverhandler.go - Interface for interaction between comms and server
//
// Copyright © 2018 Privategrity Corporation
// All rights reserved.

package mixserver

import "gitlab.com/privategrity/comms/mixmessages"

type ServerHandler interface {
	// Server Interface for the PrecompDecrypt Messages
	PrecompDecrypt(*mixmessages.PrecompDecryptMessage)
	// Server Interface for the PrecompEncrypt Messages
	PrecompEncrypt(*mixmessages.PrecompEncryptMessage)
	// Server Interface for the PrecompPermute Messages
	PrecompPermute(*mixmessages.PrecompPermuteMessage)
	// Server Interface for the RealtimeDecrypt Messages
	RealtimeDecrypt(*mixmessages.RealtimeDecryptMessage)
	// Server Interface for the RealtimeEncrypt Messages
	RealtimeEncrypt(*mixmessages.RealtimeEncryptMessage)
}
