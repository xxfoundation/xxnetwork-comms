package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	_, err := SendNewRound(SERVER_ADDRESS, &pb.InitRound{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}
