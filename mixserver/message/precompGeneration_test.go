package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompGeneration
func TestSendPrecompGeneration(t *testing.T) {
	_, err := SendPrecompGeneration(SERVER_ADDRESS, &pb.PrecompGenerationMessage{})
	if err != nil {
		t.Errorf("PrecompGeneration: Error received: %s", err)
	}
}
