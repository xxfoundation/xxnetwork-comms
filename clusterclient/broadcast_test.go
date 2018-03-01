////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package clusterclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendRealtimePermute
func TestSetPublicKey(t *testing.T) {
	_, err := SetPublicKey(SERVER_ADDRESS, &pb.PublicKeyMessage{})
	if err != nil {
		t.Errorf("PublicKeyMessage: Error received: %s", err)
	}
}

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	_, err := SendAskOnline(SERVER_ADDRESS, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

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

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	_, err := SendNewRound(SERVER_ADDRESS, &pb.InitRound{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}
