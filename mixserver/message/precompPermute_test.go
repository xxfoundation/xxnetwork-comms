package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompPermute
func TestSendPrecompPermute(t *testing.T) {
	_, err := SendPrecompPermute(SERVER_ADDRESS, &pb.PrecompPermuteMessage{})
	if err != nil {
		t.Errorf("PrecompPermute: Error received: %s", err)
	}
}
