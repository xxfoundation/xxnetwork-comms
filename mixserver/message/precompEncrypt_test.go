package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompEncrypt
func TestSendPrecompEncrypt(t *testing.T) {
	_, err := SendPrecompEncrypt(SERVER_ADDRESS, &pb.PrecompEncryptMessage{})
	if err != nil {
		t.Errorf("PrecompEncrypt: Error received: %s", err)
	}
}
