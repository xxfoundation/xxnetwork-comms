////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"context"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"io"
	"reflect"
	"testing"
)

// Creates a sender and receiver server for post phase
// unary streaming test.  The test creates a header,
// sends some slots and blocks until an ack is received
// The receive stores the slots and header data into a
// received batch to compare with expected values of batch.
func TestPhase_StreamPostPhaseSendReceive(t *testing.T) {

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := NewImplementation()
	receiverImpl.Functions.StreamPostPhase = func(server mixmessages.Node_StreamPostPhaseServer) error {
		return mockStreamPostPhase(server)
	}
	serverStreamReceiver := StartNode(servReceiverAddress, receiverImpl,
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath(), "")

	// Init server sender
	servSenderAddress := getNextServerAddress()
	serverStreamSender := StartNode(servSenderAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath(), "")

	// Get credentials and connect to node
	creds := connect.NewCredentialsFromFile(testkeys.GetNodeCertPath(),
		"*.cmix.rip")

	senderToReceiverID := MockID("sender2receiver")
	receiverToSenderID := MockID("receiver2tosender")
	// It might make more sense to call the RPC on the connection object
	// that's returned from this
	serverStreamSender.ConnectToNode(senderToReceiverID, servReceiverAddress, creds)
	serverStreamSender.ConnectToNode(receiverToSenderID, servSenderAddress, creds)

	// Reset TLS-related global variables
	defer serverStreamReceiver.Shutdown()
	defer serverStreamSender.Shutdown()

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

	streamClient, cancel, err := serverStreamSender.GetPostPhaseStreamClient(senderToReceiverID, batchInfo)

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
			t.Errorf("Failed to close and receive stream")
		}

		if ack != nil && ack.Error != "" {
			t.Errorf("Remote Server Error in ack: %s", ack.Error)
		}

	}(slots)

	// Block until stream client context is cleaned up
	<-streamClient.Context().Done()

	// Compared received slots to expected values
	for i, slot := range slots {
		if !reflect.DeepEqual(*receivedBatch.Slots[i], slot) {
			t.Errorf("Received slot %v doesn't match expected value: %v, got %v", i, slot, *receivedBatch.Slots[i])
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
	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	_ = StartNode(servReceiverAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath(), "")

	// Init server sender
	servSenderAddress := getNextServerAddress()
	serverStreamSender := StartNode(servSenderAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath(), "")

	// Get credentials and connect to node
	creds := connect.NewCredentialsFromFile(testkeys.GetNodeCertPath(),
		"*.cmix.rip")

	senderToReceiverID := MockID("sender2receiver")

	serverStreamSender.ConnectToNode(senderToReceiverID, servReceiverAddress, creds)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := serverStreamSender.getPostPhaseStream(senderToReceiverID, ctx)
	if err == nil {
		t.Errorf("Getting streaming client after canceling context should error")
	}
}

var receivedBatch mixmessages.Batch

func mockStreamPostPhase(stream mixmessages.Node_StreamPostPhaseServer) error {

	// Get header from stream
	batchInfo, err := GetPostPhaseStreamHeader(stream)
	if err != nil {
		return err
	}

	// Receive all slots and on EOF store all data
	// into a global received batch variable then
	// send ack back to client.
	var slots []*mixmessages.Slot
	for {
		slot, err := stream.Recv()
		// If we are at end of receiving
		// send ack and finish
		if err == io.EOF {
			ack := mixmessages.Ack{
				Error: "",
			}
			// Create batch using batch info header
			// and temporary slot buffer contents
			receivedBatch = mixmessages.Batch{
				Round:     batchInfo.Round,
				FromPhase: batchInfo.FromPhase,
				Slots:     slots,
			}

			err = stream.SendAndClose(&ack)

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
			Index:          i,
			MessagePayload: []byte{0x01},
			SenderID:       []byte{0x02},
		}
	}

	return slots
}
