////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, ServerAddress, nil)
	defer server.Shutdown()
	_, err := server.SendAskOnline(connID, &pb.Ping{})
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendFinishRealtime
func TestSendFinishRealtime(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	connID := MockID("node2node")
	server.ConnectToNode(connID, ServerAddress, nil)
	defer server.Shutdown()
	_, err := server.SendFinishRealtime(connID, &pb.RoundInfo{ID: 0})
	if err != nil {
		t.Errorf("FinishRealtime: Error received: %s", err)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, ServerAddress, nil)
	defer server.Shutdown()
	_, err := server.SendNewRound(connID, &pb.RoundInfo{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}

// Smoke test SendPhase
func TestSendPostPhase(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, ServerAddress, nil)
	defer server.Shutdown()
	_, err := server.SendPostPhase(connID, &pb.Batch{})
	if err != nil {
		t.Errorf("Phase: Error received: %s", err)
	}
}

// Smoke test SendPostRoundPublicKey
func TestSendPostRoundPublicKey(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, ServerAddress, nil)
	defer server.Shutdown()
	_, err := server.SendPostRoundPublicKey(connID, &pb.RoundPublicKey{})
	if err != nil {
		t.Errorf("PostRoundPublicKey: Error received: %s", err)
	}
}

// TestPostPrecompResult Smoke test
func TestSendPostPrecompResult(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode(ServerAddress, NewImplementation(), "", "", "")
	defer server.Shutdown()
	connID := MockID("connection35")
	// Connect the server to itself
	server.ConnectToNode(connID, ServerAddress, nil)
	slots := make([]*pb.Slot, 0)
	_, err := server.SendPostPrecompResult(connID, 0, slots)
	if err != nil {
		t.Errorf("PostPrecompResult: Error received: %s", err)
	}
}

func TestSendGetMeasure(t *testing.T) {
	ServerAddress := getNextServerAddress()

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()
	mockMeasure := func(msg *pb.RoundInfo) (*pb.RoundMetrics, error) {
		mockReturn := pb.RoundMetrics{
			RoundMetricJSON: "{'actual':'json'}",
		}
		return &mockReturn, nil
	}
	impl.Functions.GetMeasure = mockMeasure
	server := StartNode(ServerAddress, impl, "", "", "")
	defer server.Shutdown()

	connID := MockID("connection35")
	server.ConnectToNode(connID, ServerAddress, nil)
	ri := pb.RoundInfo{
		ID: uint64(3),
	}
	_, err := server.SendGetMeasure(connID, &ri)
	if err != nil {
		t.Errorf("SendGetMeasure: Error received: %s", err)
	}
}

func TestSendGetMeasureError(t *testing.T) {
	ServerAddress := getNextServerAddress()

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()

	mockMeasureError := func(msg *pb.RoundInfo) (*pb.RoundMetrics, error) {
		return nil, errors.New("Test error")
	}
	impl.Functions.GetMeasure = mockMeasureError
	server := StartNode(ServerAddress, impl, "", "", "")
	defer server.Shutdown()

	connID := MockID("connection35")
	server.ConnectToNode(connID, ServerAddress, nil)
	ri := pb.RoundInfo{
		ID: uint64(3),
	}
	_, err := server.SendGetMeasure(connID, &ri)
	if err == nil {
		t.Error("Did not receive error response")
	}
}
