////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

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

// Smoke test SendRealtimeEncrypt
func TestSendRealtimeEncrypt(t *testing.T) {
	_, err := SendRealtimeEncrypt(SERVER_ADDRESS, &pb.RealtimeEncryptMessage{})
	if err != nil {
		t.Errorf("RealtimeEncrypt: Error received: %s", err)
	}
}

// Smoke test SendRealtimeDecrypt
func TestSendRealtimeDecrypt(t *testing.T) {
	_, err := SendRealtimeDecrypt(SERVER_ADDRESS, &pb.RealtimeDecryptMessage{})
	if err != nil {
		t.Errorf("RealtimeDecrypt: Error received: %s", err)
	}
}
