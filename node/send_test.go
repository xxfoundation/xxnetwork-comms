///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"testing"
)

// Smoke test SendAskOnline
func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testID, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
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
	testID := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testID, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()
	mockMeasure := func(message *pb.RoundInfo, auth *connect.Auth) (*pb.RoundMetrics, error) {
		mockReturn := pb.RoundMetrics{
			RoundMetricJSON: "{'actual':'json'}",
		}
		return &mockReturn, nil
	}
	impl.Functions.GetMeasure = mockMeasure
	server := StartNode(testId, ServerAddress, impl, nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)

	// GRPC complains if this doesn't return something nice, so I mocked it
	impl := NewImplementation()

	mockMeasureError := func(message *pb.RoundInfo, auth *connect.Auth) (*pb.RoundMetrics, error) {
		return nil, errors.New("Test error")
	}
	impl.Functions.GetMeasure = mockMeasureError
	server := StartNode(testId, ServerAddress, impl, nil, nil)
	defer server.Shutdown()

	ri := pb.RoundInfo{
		ID: uint64(3),
	}
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
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
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, impl, nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	any, err := ptypes.MarshalAny(&messages.Ack{})
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

func TestSendRoundError(t *testing.T) {
	ServerAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = server.SendRoundError(host, &pb.RoundError{})
	if err != nil {
		t.Errorf("NewRound: Error received: %s", err)
	}

}
