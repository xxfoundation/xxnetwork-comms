package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompDecrypt
func TestSendPrecompDecrypt(t *testing.T) {
	_, err := SendPrecompDecrypt(SERVER_ADDRESS, &pb.PrecompDecryptMessage{})
	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
}
