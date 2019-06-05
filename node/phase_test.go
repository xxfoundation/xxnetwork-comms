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
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestPhase_StreamPostPhase(t *testing.T) {

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	serverStreamReceiver := StartNode(servReceiverAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())

	// Init server sender
	servSenderAddress := getNextServerAddress()
	serverStreamSender := StartNode(servSenderAddress, NewImplementation(),
		testkeys.GetNodeCertPath(), testkeys.GetNodeKeyPath())

	// Get credentials and connenct to node
	creds := connect.NewCredentialsFromFile(testkeys.GetNodeCertPath(),
		"*.cmix.rip")

	connectionID := MockID("server2toserver")
	// It might make more sense to call the RPC on the connection object
	// that's returned from this
	serverStreamSender.ConnectToNode(connectionID, servReceiverAddress, creds)
	// Reset TLS-related global variables
	defer serverStreamReceiver.Shutdown()
	defer serverStreamSender.Shutdown()

	// Create streaming context so you can close stream later
	ctx, cancel := connect.StreamingContext()

	// Create round info to be used in batch info
	roundInfo := mixmessages.RoundInfo{
		ID: uint64(10),
	}

	// Create header which contains the batch info
	batchInfo := mixmessages.BatchInfo{
		Round:    &roundInfo,
		ForPhase: 3,
	}

	// Create a new context with some metadata
	ctx = metadata.AppendToOutgoingContext(ctx,
		"BatchInfo", batchInfo.String())

	// Get stream client for post phase
	_, _ = serverStreamSender.GetPostPhaseStream(connectionID, ctx)

	//slots := createSlots(3)
	//
	//go func(slots []mixmessages.Slot) {
	//	for i, slot := range slots {
	//		err = streamClient.Send(&slot)
	//		if err != nil {
	//			t.Errorf("Unable to send slot %v %v", i, err)
	//		}
	//	}
	//	//_, err := streamClient.CloseAndRecv()
	//	//// todo: check ack error?
	//	//if err != nil {
	//	//	t.Errorf("Unable to close stream%v", err)
	//	//}
	//
	//}(slots)

	// TODO: Check if server receiver gets all 3 slots and header

	// Wait until
	//<-ctx.Done()
	//<-streamClient.Context().Done()
	//_, err = streamClient.CloseAndRecv()
	//// todo: check ack error?
	//if err != nil {
	//	t.Errorf("Unable to close stream%v", err)
	//}
	//
	//err = ctx.Err()
	//if err != nil {
	//	t.Errorf("Error when sending slots")
	//}

	cancel()
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
