////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"bytes"
	"context"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	jww.SetStdoutThreshold(jww.LevelTrace)
	connect.TestingOnlyDisableTLS = true
	os.Exit(m.Run())
}

func TestDownloadBatch(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	nodeID := id.NewIdFromString("test", id.Node, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), certData,
		keyData, gossip.DefaultManagerFlags())

	batchSize := uint32(3)
	slots := make([]*mixmessages.Slot, batchSize)
	for i := uint32(0); i < batchSize; i++ {
		slots[i] = &mixmessages.Slot{
			Index:    i,
			PayloadA: []byte{0x01},
			SenderID: []byte{0x02},
		}
	}
	receiverImpl := node.NewImplementation()
	receiverImpl.Functions.DownloadMixedBatch = func(stream mixmessages.Node_DownloadMixedBatchServer,
		batchInfo *mixmessages.BatchReady, auth *connect.Auth) error {
		return mockStreamMixedBatch(stream, slots, t)
	}

	server := node.StartNode(nodeID, ServerAddress, 0, receiverImpl,
		certData, keyData)
	defer gateway.Shutdown()
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(nodeID, ServerAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	received, err := gateway.DownloadMixedBatch(&mixmessages.BatchReady{RoundId: 4}, host)
	if err != nil {
		t.Fatalf("DownloadMixedBatch: Error received: %v", err)
	}

	if len(slots) != len(received) {
		t.Errorf("Did not receive expected number of slots."+
			"\n\tExpected: %v"+
			"\n\tReceived: %v", len(slots), len(received))
	}

}

func mockStreamMixedBatch(server mixmessages.Node_DownloadMixedBatchServer,
	slots []*mixmessages.Slot, t *testing.T) error {
	// Receive all slots and on EOF store all data
	// into a global received batch variable then
	// send ack back to client.
	for i, slot := range slots {
		err := server.Send(slot)
		if err != nil {
			t.Errorf("Unable to send slot %v %v", i, err)
		}
	}

	return nil
}

// Happy path
func TestComms_StreamUnmixedBatch(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := node.NewImplementation()
	receiverImpl.Functions.UploadUnmixedBatch = func(server mixmessages.Node_UploadUnmixedBatchServer, auth *connect.Auth) error {
		return mockStreamUnmixedBatch(server)
	}

	testID := id.NewIdFromString("test", id.Generic, t)
	serverStreamReceiver := node.StartNode(testID, servReceiverAddress, 0, receiverImpl,
		certData, keyData)

	// Init sender
	senderAddress := getNextServerAddress()
	gwStreamSender := StartGateway(testID, senderAddress,
		NewImplementation(), nil, nil, gossip.DefaultManagerFlags())

	// Reset TLS-related global variables
	defer serverStreamReceiver.Shutdown()
	defer gwStreamSender.Shutdown()

	// Create header
	roundId := uint64(10)
	roundInfo := mixmessages.RoundInfo{
		ID: roundId,
	}
	fromPhase := int32(3)
	batchSize := uint32(3)
	batchInfo := mixmessages.BatchInfo{
		Round:     &roundInfo,
		FromPhase: fromPhase,
		BatchSize: batchSize,
	}

	// Init host/manager
	manager := connect.NewManagerTesting(t)
	testId := id.NewIdFromString("test", id.Generic, t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, servReceiverAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	mockBatch := &mixmessages.Batch{
		Round:     &roundInfo,
		FromPhase: fromPhase,
	}

	for i := uint32(0); i < batchSize; i++ {
		mockBatch.Slots = append(mockBatch.Slots,
			&mixmessages.Slot{
				Index:    i,
				PayloadA: []byte{byte(i)},
			})
	}

	err = gwStreamSender.UploadUnmixedBatch(host, batchInfo, mockBatch)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
}

// Creates a sender and receiver server for post phase
// unary streaming test.  The test creates a header,
// sends some slots and blocks until an ack is received
// The receive stores the slots and header data into a
// received batch to compare with expected values of batch.
func TestPhase_StreamPostPhaseSendReceive(t *testing.T) {

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := node.NewImplementation()
	receiverImpl.Functions.UploadUnmixedBatch = func(server mixmessages.Node_UploadUnmixedBatchServer, auth *connect.Auth) error {
		return mockStreamUnmixedBatch(server)
	}

	testID := id.NewIdFromString("test", id.Generic, t)
	serverStreamReceiver := node.StartNode(testID, servReceiverAddress, 0, receiverImpl,
		certData, keyData)

	// Init sender
	senderAddress := getNextServerAddress()
	gwStreamSender := StartGateway(testID, senderAddress,
		NewImplementation(), nil, nil, gossip.DefaultManagerFlags())

	// Reset TLS-related global variables
	defer serverStreamReceiver.Shutdown()
	defer gwStreamSender.Shutdown()

	// Create header
	roundId := uint64(10)
	roundInfo := mixmessages.RoundInfo{
		ID: roundId,
	}
	fromPhase := int32(3)
	batchSize := uint32(3)
	batchInfo := mixmessages.BatchInfo{
		Round:     &roundInfo,
		FromPhase: fromPhase,
		BatchSize: batchSize,
	}

	// Init host/manager
	manager := connect.NewManagerTesting(t)
	testId := id.NewIdFromString("test", id.Generic, t)
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testId, servReceiverAddress, certData, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	streamClient, cancel, err := gwStreamSender.getUnmixedBatchStreamClient(
		host, batchInfo)

	if err != nil {
		t.Errorf("Unable to get streaming client %v", err)
	}

	// Generate indexed slots
	slots := createSlots(batchSize)

	// The server will send the slot messages and wait for ack
	// and close on receiving it.
	go func(slots []mixmessages.Slot) {
		for i, slot := range slots {
			err := streamClient.Send(&slot)
			if err != nil {
				t.Errorf("Unable to send slot %v %v", i, err)
			}
		}
		ack, err := streamClient.CloseAndRecv()

		if err != nil {
			t.Errorf("Failed to close and receive stream: %+v", err)
		}

		if ack != nil && ack.Error != "" {
			t.Errorf("Remote Server Error in ack: %s", ack.Error)
		}

	}(slots)

	// Block until stream client context is cleaned up
	<-streamClient.Context().Done()

	// Compared received slots to expected values
	for i, slot := range slots {
		if !slotCmp(*receivedBatch.Slots[i], slot) {
			t.Errorf("Received slot %v doesn't match expected value: %v, got %v", i, slot, *receivedBatch.Slots[i])
			//t.Errorf("Received slot %v doesn't match expected\n\texpected: %#v\n\treceived: %#v", i, slot, *receivedBatch.Slots[i])
		}
	}

	// Compare for phase
	if fromPhase != receivedBatch.FromPhase {
		t.Errorf("FromPhase received %v doesn't match expected %v", receivedBatch.FromPhase, fromPhase)
	}

	// Compare round info
	if !reflect.DeepEqual(roundInfo, *receivedBatch.Round) {
		t.Errorf("Round info received %v doesn't match expected %v", *receivedBatch.Round, roundInfo)
	}

	// Clean up sender context
	cancel()
}

// getPostPhaseStream should error when context canceled before call
func TestGetPostPhaseStream_ErrorsWhenContextCanceled(t *testing.T) {
	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Node, t)
	_ = node.StartNode(testID, servReceiverAddress, 0, node.NewImplementation(), certData, keyData)

	// Init server sender
	gwSenderAddress := getNextServerAddress()
	gwStreamSender := StartGateway(testID, gwSenderAddress,
		NewImplementation(), nil, nil, gossip.DefaultManagerFlags())

	// Get credentials and connect to node
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Init host/manager
	manager := connect.NewManagerTesting(t)
	host, err := manager.AddHost(testID, servReceiverAddress, certData,
		connect.GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = gwStreamSender.getUnmixedBatchStream(host, ctx)
	if err == nil {
		t.Errorf("Getting streaming client after canceling context should error")
	}
}

var receivedBatch mixmessages.Batch

func mockStreamUnmixedBatch(server mixmessages.Node_UploadUnmixedBatchServer) error {
	// Get header from stream
	batchInfo, err := node.GetUnmixedBatchStreamHeader(server)
	if err != nil {
		return err
	}

	// Receive all slots and on EOF store all data
	// into a global received batch variable then
	// send ack back to client.
	var slots []*mixmessages.Slot
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
			receivedBatch = mixmessages.Batch{
				Round:     batchInfo.Round,
				FromPhase: batchInfo.FromPhase,
				Slots:     slots,
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

// createSlots is a helper function to generate slot
// messages uses for testing
func createSlots(numSlots uint32) []mixmessages.Slot {

	slots := make([]mixmessages.Slot, numSlots)

	for i := uint32(0); i < numSlots; i++ {
		slots[i] = mixmessages.Slot{
			Index:    i,
			PayloadA: []byte{0x01},
			SenderID: []byte{0x02},
		}
	}

	return slots
}

// slotCmp compares the given slots to each other and returns true if the public
// fields are the same.
func slotCmp(slotA, slotB mixmessages.Slot) bool {
	if slotA.GetIndex() != slotB.GetIndex() {
		return false
	} else if !bytes.Equal(slotA.GetEncryptedPayloadAKeys(), slotB.GetEncryptedPayloadAKeys()) {
		return false
	} else if !bytes.Equal(slotA.GetEncryptedPayloadBKeys(), slotB.GetEncryptedPayloadBKeys()) {
		return false
	} else if !bytes.Equal(slotA.GetPartialPayloadACypherText(), slotB.GetPartialPayloadACypherText()) {
		return false
	} else if !bytes.Equal(slotA.GetPartialPayloadBCypherText(), slotB.GetPartialPayloadBCypherText()) {
		return false
	} else if !bytes.Equal(slotA.GetPartialRoundPublicCypherKey(), slotB.GetPartialRoundPublicCypherKey()) {
		return false
	} else if !bytes.Equal(slotA.GetSenderID(), slotB.GetSenderID()) {
		return false
	} else if !bytes.Equal(slotA.GetPayloadA(), slotB.GetPayloadA()) {
		return false
	} else if !bytes.Equal(slotA.GetPayloadB(), slotB.GetPayloadB()) {
		return false
	} else if !bytes.Equal(slotA.GetSalt(), slotB.GetSalt()) {
		return false
	}

	return true
}
