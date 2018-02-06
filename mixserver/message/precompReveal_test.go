package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompReveal
func TestSendPrecompReveal(t *testing.T) {
	_, err := SendPrecompReveal(SERVER_ADDRESS, &pb.PrecompRevealMessage{})
	if err != nil {
		t.Errorf("PrecompReveal: Error received: %s", err)
	}
}
