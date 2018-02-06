package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRealtimePermute
func TestSendRealtimePermute(t *testing.T) {
	_, err := SendRealtimePermute(SERVER_ADDRESS, &pb.RealtimePermuteMessage{})
	if err != nil {
		t.Errorf("RealtimePermute: Error received: %s", err)
	}
}
