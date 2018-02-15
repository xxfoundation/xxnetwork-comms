package message

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRequestMessage
func TestSendRequestMessage(t *testing.T) {
	_, err := SendRequestMessage(SERVER_ADDRESS, &pb.RequestMessage{})
	if err != nil {
		t.Errorf("RequestMessage: Error received: %s", err)
	}
}

// Smoke test SendMessagetoSender
func TestSendMessageToServer(t *testing.T) {
	_, err := SendMessageToServer(SERVER_ADDRESS, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("SendMessageToServer: Error received: %s", err)
	}
}
