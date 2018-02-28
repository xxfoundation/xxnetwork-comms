package mixclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendClientPoll
func TestSendClientPoll(t *testing.T) {
	_, err := SendClientPoll(SERVER_ADDRESS, &pb.ClientPollMessage{})
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
