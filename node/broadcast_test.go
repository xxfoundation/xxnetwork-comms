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
func TestSetPublicKey(t *testing.T) {
	_, err := SetPublicKey(ServerAddress, &pb.PublicKeyMessage{})
	if err != nil {
		t.Errorf("PublicKeyMessage: Error received: %s", err)
	}
}

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	_, err := SendAskOnline(ServerAddress, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendRoundtripPing
func TestSendRoundtripPing(t *testing.T) {
	_, err := SendRoundtripPing(ServerAddress, &pb.TimePing{})
	if err != nil {
		t.Errorf("RoundtripPing: Error received: %s", err)
	}
}

// Smoke test SendServerMetrics
func TestSendServerMetrics(t *testing.T) {
	_, err := SendServerMetrics(ServerAddress, &pb.ServerMetricsMessage{})
	if err != nil {
		t.Errorf("ServerMetrics: Error received: %s", err)
	}
}

// Smoke test SendNetworkError
func TestSendNetworkError(t *testing.T) {
	r, err := SendNetworkError(ServerAddress, &pb.ErrorMessage{Message: "Hello, world!"})

	if err != nil {
		t.Errorf("PrecompDecrypt: Error received: %s", err)
	}
	if r.MsgLen != 13 {
		t.Errorf("NetworkError: Expected len of %v, got %v", 13, r)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	_, err := SendNewRound(ServerAddress, &pb.InitRound{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}

// Smoke test SendUserUpsert
func TestSendUserUpsert(t *testing.T) {
	_, err := SendUserUpsert(ServerAddress, &pb.UpsertUserMessage{})
	if err != nil {
		t.Errorf("UserUpsert: Error received: %s", err)
	}
}
