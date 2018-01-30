package mixserver

import "gitlab.com/privategrity/comms/mixmessages"

type ServerHandler interface {
	// Server Interface for the PrecompDecrypt Messages
	precompDecrypt(mixmessages.PrecompDecryptMessage)
}
