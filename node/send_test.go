////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"io"
	"testing"
)

// Smoke test SendAskOnline
//todo: fix and re enable
/*func TestSendAskOnline(t *testing.T) {
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testID, ServerAddress, 0, NewImplementation(), nil, nil)
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
}*/

func TestComms_StreamPrecompTestBatch(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Construct sender
	servSenderAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Node, t)
	senderServer := StartNode(testID, servSenderAddress, 0, NewImplementation(), nil, nil)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := NewImplementation()
	receiverImpl.Functions.PrecompTestBatch = func(server pb.Node_PrecompTestBatchServer, info *pb.RoundInfo,
		auth *connect.Auth) error {
		return mockStreamPrecompTestBatch(server)
	}
	serverStreamReceiver := StartNode(testID, servReceiverAddress, 0, receiverImpl,
		certData, keyData)

	defer senderServer.Shutdown()
	defer serverStreamReceiver.Shutdown()

	// Init host/manager
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, servReceiverAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	mockBatch := &pb.CompletedBatch{}

	err = senderServer.StreamPrecompTestBatch(host, &pb.RoundInfo{ID: 0}, mockBatch)
	if err != nil {
		t.Errorf("StreamPrecompTestBatch: Error received: %s", err)
	}
}

func mockStreamPrecompTestBatch(server pb.Node_PrecompTestBatchServer) error {
	// Get header from stream
	roundInfo, err := GetPrecompTestBatchStreamHeader(server)
	if err != nil {
		return err
	}

	// Receive all slots and on EOF store all data
	// into a global received batch variable then
	// send ack back to client.
	var slots []*pb.Slot
	for {
		slot, err := server.Recv()
		// If we are at end of receiving
		// send ack and finish
		if err == io.EOF {
			ack := messages.Ack{
				Error: "",
			}
			// Create batch using batch info header
			// and temporary slot buffer contents
			completedBatchRealtime = pb.CompletedBatch{
				Slots:   slots,
				RoundID: roundInfo.ID,
			}

			err = server.SendAndClose(&ack)

			return err
		}

		// If we have another error, return err
		if err != nil {
			return err
		}

		// Store slot received into temporary buffer
		slots = append(slots, slot)
	}
}

// Smoke test SendFinishRealtime
func TestSendFinishRealtime(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Construct sender
	servSenderAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Node, t)
	senderServer := StartNode(testID, servSenderAddress, 0, NewImplementation(), nil, nil)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := NewImplementation()
	receiverImpl.Functions.FinishRealtime = func(roundInfo *pb.RoundInfo, server pb.Node_FinishRealtimeServer, auth *connect.Auth) error {
		return mockStreamFinishRealtime(server)
	}
	serverStreamReceiver := StartNode(testID, servReceiverAddress, 0, receiverImpl,
		certData, keyData)

	defer senderServer.Shutdown()
	defer serverStreamReceiver.Shutdown()

	// Init host/manager
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, servReceiverAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	mockBatch := &pb.CompletedBatch{}

	_, err = senderServer.SendFinishRealtime(host, &pb.RoundInfo{ID: 0}, mockBatch)
	if err != nil {
		t.Errorf("FinishRealtime: Error received: %s", err)
	}
}

var completedBatchRealtime pb.CompletedBatch

func mockStreamFinishRealtime(server pb.Node_FinishRealtimeServer) error {
	// Get header from stream
	roundInfo, err := GetFinishRealtimeStreamHeader(server)
	if err != nil {
		return err
	}

	// Receive all slots and on EOF store all data
	// into a global received batch variable then
	// send ack back to client.
	var slots []*pb.Slot
	for {
		slot, err := server.Recv()
		// If we are at end of receiving
		// send ack and finish
		if err == io.EOF {
			ack := messages.Ack{
				Error: "",
			}
			// Create batch using batch info header
			// and temporary slot buffer contents
			completedBatchRealtime = pb.CompletedBatch{
				Slots:   slots,
				RoundID: roundInfo.ID,
			}

			err = server.SendAndClose(&ack)

			return err
		}

		// If we have another error, return err
		if err != nil {
			return err
		}

		// Store slot received into temporary buffer
		slots = append(slots, slot)
	}

}

// Smoke test SendNewRound
func TestSendNewRound(t *testing.T) {
	ServerAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, 0, NewImplementation(), nil, nil)
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
	server := StartNode(testId, ServerAddress, 0, NewImplementation(), nil, nil)
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

// TestPostPrecompResult Smoke test
func TestSendPostPrecompResult(t *testing.T) {
	ServerAddress := getNextServerAddress()
	testId := id.NewIdFromString("test", id.Node, t)
	server := StartNode(testId, ServerAddress, 0, NewImplementation(), nil, nil)
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}
	slots := make([]*pb.Slot, 0)
	_, err = server.SendPostPrecompResult(host, 0, uint32(len(slots)))
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
	server := StartNode(testId, ServerAddress, 0, impl, nil, nil)
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
	server := StartNode(testId, ServerAddress, 0, impl, nil, nil)
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
	server := StartNode(testId, ServerAddress, 0, impl, nil, nil)
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
	server := StartNode(testId, ServerAddress, 0, NewImplementation(), nil, nil)
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
