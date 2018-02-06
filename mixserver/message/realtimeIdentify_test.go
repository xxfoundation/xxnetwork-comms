package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRealtimeIdentify
func TestSendRealtimeIdentify(t *testing.T) {
	_, err := SendRealtimeIdentify(SERVER_ADDRESS, &pb.RealtimeIdentifyMessage{})
	if err != nil {
		t.Errorf("RealtimeIdentify: Error received: %s", err)
	}
}
