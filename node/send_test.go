////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendAskOnline(connID, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendRoundtripPing
func TestSendRoundtripPing(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendRoundtripPing(connID, &pb.TimePing{})
	if err != nil {
		t.Errorf("RoundtripPing: Error received: %s", err)
	}
}

// Smoke test SendFinishRealtime
func TestSendFinishRealtime(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("node2node")
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendFinishRealtime(connID)
	if err != nil {
		t.Errorf("FinishRealtime: Error received: %s", err)
	}
}

// Smoke test SendServerMetrics
func TestSendServerMetrics(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendServerMetrics(connID, &pb.ServerMetrics{})
	if err != nil {
		t.Errorf("ServerMetrics: Error received: %s", err)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendNewRound(connID, &pb.RoundInfo{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}

// Smoke test SendPhase
func TestSendPostPhase(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendPostPhase(connID, &pb.Batch{})
	if err != nil {
		t.Errorf("Phase: Error received: %s", err)
	}
}

// Smoke test SendPostRoundPublicKey
func TestSendPostRoundPublicKey(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	defer server.Shutdown()
	_, err := server.SendPostRoundPublicKey(connID, &pb.RoundPublicKey{})
	if err != nil {
		t.Errorf("PostRoundPublicKey: Error received: %s", err)
	}
}

// TestPostPrecompResult Smoke test
func TestSendPostPrecompResult(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartServer(ServerAddress, NewImplementation(), "", "")
	defer server.Shutdown()
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, &connect.ConnectionInfo{
		Address: ServerAddress,
	})
	slots := make([]*pb.Slot, 0)
	_, err := server.SendPostPrecompResult(connID, 0, slots)
	if err != nil {
		t.Errorf("PostPrecompResult: Error received: %s", err)
	}
}
