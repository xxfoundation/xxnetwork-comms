package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompDecrypt
func TestSendPrecompDecrypt(t *testing.T) {
	addr := "localhost:5555"
	_, err := SendPrecompDecrypt(addr, &pb.PrecompDecryptMessage{})
	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
}
