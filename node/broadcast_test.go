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

// Smoke test SendRealtimePermute
func TestSetPublicKey(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SetPublicKey(ServerAddress, &pb.PublicKeyMessage{})
	if err != nil {
		t.Errorf("PublicKeyMessage: Error received: %s", err)
	}
}

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendAskOnline(ServerAddress, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendRoundtripPing
func TestSendRoundtripPing(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendRoundtripPing(ServerAddress, &pb.TimePing{})
	if err != nil {
		t.Errorf("RoundtripPing: Error received: %s", err)
	}
}

// Smoke test SendServerMetrics
func TestSendServerMetrics(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendServerMetrics(ServerAddress, &pb.ServerMetricsMessage{})
	if err != nil {
		t.Errorf("ServerMetrics: Error received: %s", err)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ShutDown := StartServer(ServerAddress, NewImplementation(), "", "")
	defer ShutDown()
	_, err := SendNewRound(ServerAddress, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}
