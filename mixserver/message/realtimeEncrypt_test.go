package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRealtimeEncrypt
func TestSendRealtimeEncrypt(t *testing.T) {
	_, err := SendRealtimeEncrypt(SERVER_ADDRESS, &pb.RealtimeEncryptMessage{})
	if err != nil {
		t.Errorf("RealtimeEncrypt: Error received: %s", err)
	}
}
