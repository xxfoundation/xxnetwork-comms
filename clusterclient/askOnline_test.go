package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	_, err := SendAskOnline(SERVER_ADDRESS, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}
