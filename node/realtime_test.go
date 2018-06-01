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
	ShutDown := StartServer(ServerAddress, NewImplementation())
	defer ShutDown()
	_, err := SendRealtimePermute(ServerAddress, &pb.RealtimePermuteMessage{})
	if err != nil {
		t.Errorf("RealtimePermute: Error received: %s", err)
	}
}

// Smoke test SendRealtimeEncrypt
func TestSendRealtimeEncrypt(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation())
	defer ShutDown()
	_, err := SendRealtimeEncrypt(ServerAddress, &pb.RealtimeEncryptMessage{})
	if err != nil {
		t.Errorf("RealtimeEncrypt: Error received: %s", err)
	}
}

// Smoke test SendRealtimeDecrypt
func TestSendRealtimeDecrypt(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation())
	defer ShutDown()
	_, err := SendRealtimeDecrypt(ServerAddress, &pb.RealtimeDecryptMessage{})
	if err != nil {
		t.Errorf("RealtimeDecrypt: Error received: %s", err)
	}
}
