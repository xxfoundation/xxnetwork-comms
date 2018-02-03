package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRealtimeDecrypt
func TestSendRealtimeDecrypt(t *testing.T) {
	_, err := SendRealtimeDecrypt(SERVER_ADDRESS, &pb.RealtimeDecryptMessage{})
	if err != nil {
		t.Errorf("RealtimeDecrypt: Error received: %s", err)
	}
}
