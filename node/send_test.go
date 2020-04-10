////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendAskOnline(host)
	if err != nil {
		t.Errorf("AskOnline: Error received: %s", err)
	}
}

// Smoke test SendFinishRealtime
func TestSendFinishRealtime(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendFinishRealtime(host, &pb.RoundInfo{ID: 0})
	if err != nil {
		t.Errorf("FinishRealtime: Error received: %s", err)
	}
}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendNewRound(host, &pb.RoundInfo{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}
}

// Smoke test SendPhase
func TestSendPostPhase(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendPostPhase(host, &pb.Batch{})
	if err != nil {
		t.Errorf("Phase: Error received: %s", err)
	}
}

// Smoke test SendPostRoundPublicKey
func TestSendPostRoundPublicKey(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendPostRoundPublicKey(host, &pb.RoundPublicKey{})
	if err != nil {
		t.Errorf("PostRoundPublicKey: Error received: %s", err)
	}
}

// TestPostPrecompResult Smoke test
func TestSendPostPrecompResult(t *testing.T) {
	ServerAddress := getNextServerAddress()
	server := StartNode("test", ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	slots := make([]*pb.Slot, 0)
	_, err = server.SendPostPrecompResult(host, 0, slots)
	if err != nil {
		t.Errorf("PostPrecompResult: Error received: %s", err)
	}
}

func TestSendGetMeasure(t *testing.T) {
	ServerAddress := getNextServerAddress()

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()
	mockMeasure := func(message *pb.RoundInfo, auth *connect.Auth) (*pb.RoundMetrics, error) {
		mockReturn := pb.RoundMetrics{
			RoundMetricJSON: "{'actual':'json'}",
		}
		return &mockReturn, nil
	}
	impl.Functions.GetMeasure = mockMeasure
	server := StartNode("test", ServerAddress, impl, nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	ri := pb.RoundInfo{
		ID: uint64(3),
	}
	_, err = server.SendGetMeasure(host, &ri)
	if err != nil {
		t.Errorf("SendGetMeasure: Error received: %s", err)
	}
}

func TestSendGetMeasureError(t *testing.T) {
	ServerAddress := getNextServerAddress()

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()

	mockMeasureError := func(message *pb.RoundInfo, auth *connect.Auth) (*pb.RoundMetrics, error) {
		return nil, errors.New("Test error")
	}
	impl.Functions.GetMeasure = mockMeasureError
	server := StartNode("test", ServerAddress, impl, nil, nil)
	defer server.Shutdown()

	ri := pb.RoundInfo{
		ID: uint64(3),
	}
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendGetMeasure(host, &ri)
	if err == nil {
		t.Error("Did not receive error response")
	}
}

func TestRoundTripPing(t *testing.T) {
	ServerAddress := getNextServerAddress()
	impl := NewImplementation()
	server := StartNode("test", ServerAddress, impl, nil, nil)
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	any, err := ptypes.MarshalAny(&pb.Ack{})
	if err != nil {
		t.Errorf("SendRoundTripPing: failed attempting to marshall any type: %+v", err)
	}

	rtPing := &pb.RoundTripPing{
		Round: &pb.RoundInfo{
			ID: uint64(1),
		},
		Payload: any,
	}

	_, err = server.RoundTripPing(host, rtPing)
	if err != nil {
		t.Errorf("Received error from RoundTripPing: %+v", err)
	}
}
