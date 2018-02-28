package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendPrecompShare
func TestSendPrecompShare(t *testing.T) {
	_, err := SendPrecompShare(SERVER_ADDRESS, &pb.PrecompShareMessage{})
	if err != nil {
		t.Errorf("PrecompShare: Error received: %s", err)
	}
}
