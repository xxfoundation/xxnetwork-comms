////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"io"
	"testing"
)

func TestPhase_StreamPostPhase(t *testing.T) {

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := NewImplementation()
	receiverImpl.Functions.StreamPostPhase = func(server mixmessages.Node_StreamPostPhaseServer) error {
		return mockStreamPostPhase(server)
	}
	serverStreamReceiver := StartNode(servReceiverAddress, receiverImpl,
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())

	// Init server sender
	servSenderAddress := getNextServerAddress()
	serverStreamSender := StartNode(servSenderAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())

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

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	//// Create round info to be used in batch info
	//roundInfo := mixmessages.RoundInfo{
	//	ID: uint64(10),
	//}
	//
	//// Create header which contains the batch info
	//batchInfo := mixmessages.BatchInfo{
	//	Round:    &roundInfo,
	//	ForPhase: 3,
	//}
	//
	//// Create a new context with some metadata
	//ctx = metadata.AppendToOutgoingContext(ctx,
	//	"BatchInfo", batchInfo.String())

	// Get stream client for post phase
	streamClient, err := serverStreamSender.GetPostPhaseStream(senderToReceiverID, ctx)

	if err != nil {
		t.Errorf("Unable to get streamling clinet %v", err)
	}

	slots := createSlots(3)

	// The server will send the slot messages and wait for ack
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
			t.Errorf("Remote Server Error: %s", ack.Error)
		}

	}(slots)

	// Block until stream client context is cleaned up
	<-streamClient.Context().Done()

	// Clean up sender context
	cancel()
}

var receivedBatch mixmessages.Batch

func mockStreamPostPhase(server mixmessages.Node_StreamPostPhaseServer) error {
	receivedBatch = mixmessages.Batch{
		Slots: make([]*mixmessages.Slot, 3),
	}

	// Send a chunk per received slot until EOF
	index := uint32(0)
	for {
		slot, err := server.Recv()
		// If we are at end of receiving
		// send ack and finish
		if err == io.EOF {
			ack := mixmessages.Ack{}
			err = server.SendAndClose(&ack)
			return err
		}

		// If we have another error, return err
		if err != nil {
			return err
		}

		// Store slot received into received batch
		receivedBatch.Slots[index] = slot
		index++
	}
}

func createSlots(numSlots int) []mixmessages.Slot {

	slots := make([]mixmessages.Slot, numSlots)

	for i := 0; i < numSlots; i++ {
		slots[i] = mixmessages.Slot{
			MessagePayload: []byte{0x01},
			SenderID:       []byte{0x02},
		}
	}

	return slots
}
