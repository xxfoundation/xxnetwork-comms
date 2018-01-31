// serverhandler.go - Interface for interaction between comms and server
//
// Copyright © 2018 Privategrity Corporation
// All rights reserved.

package mixserver

import "gitlab.com/privategrity/comms/mixmessages"

type ServerHandler interface {
	// Server Interface for the PrecompDecrypt Messages
	PrecompDecrypt(mixmessages.PrecompDecryptMessage)
}
