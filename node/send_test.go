////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendAskOnline(ServerAddress, "", &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendRoundtripPing
func TestSendRoundtripPing(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendRoundtripPing(ServerAddress, "", &pb.TimePing{})
	if err != nil {
		t.Errorf("RoundtripPing: Error received: %s", err)
	}
}

// Smoke test SendServerMetrics
func TestSendServerMetrics(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendServerMetrics(ServerAddress, "", &pb.ServerMetrics{})
	if err != nil {
		t.Errorf("ServerMetrics: Error received: %s", err)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendNewRound(ServerAddress, "", &pb.Batch{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}

// Smoke test SendPhase
func TestSendPostPhase(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPostPhase(ServerAddress, "", &pb.Batch{})
	if err != nil {
		t.Errorf("Phase: Error received: %s", err)
	}
}

// Smoke test SendPostRoundPublicKey
func TestSendPostRoundPublicKey(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendPostRoundPublicKey(ServerAddress, "", &pb.RoundPublicKey{})
	if err != nil {
		t.Errorf("PostRoundPublicKey: Error received: %s", err)
	}
}

// TestFinishPrecomputation Smoke test
func TestSendFinishPrecomputation(t *testing.T) {
	ServerAddress := getNextServerAddress()
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	slots := make([]*pb.Slot, 0)
	_, err := SendFinishPrecomputation(ServerAddress, "", 0, slots)
	if err != nil {
		t.Errorf("FinishPrecomputation: Error received: %s", err)
	}
}
