package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	addr := "localhost:5555"
	_, err := SendAskOnline(addr, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}
