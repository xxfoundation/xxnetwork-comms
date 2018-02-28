package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendNetworkError
func TestSendNetworkError(t *testing.T) {
	r, err := SendNetworkError(SERVER_ADDRESS, &pb.ErrorMessage{Message: "Hello, world!"})

	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
	if r.MsgLen != 13 {
		t.Errorf("NetworkError: Expected len of %v, got %v", 13, r)
	}
}
