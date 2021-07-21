package node

import (
	"gitlab.com/elixxir/comms/gateway"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/primitives/id"
	"io"
	"testing"
)

// Smoke test DownloadMixedBatch
func TestDownloadMixedBatch(t *testing.T) {

	keyPath := testkeys.GetNodeKeyPath()
	keyData := testkeys.LoadFromPath(keyPath)
	certPath := testkeys.GetNodeCertPath()
	certData := testkeys.LoadFromPath(certPath)

	// Init server receiver
	servReceiverAddress := getNextServerAddress()
	receiverImpl := gateway.NewImplementation()
	var receivedSlots []*mixmessages.Slot
	var err error
	receiverImpl.Functions.DownloadMixedBatch = func(server mixmessages.Gateway_DownloadMixedBatchServer, auth *connect.Auth) error {
		receivedSlots, err =  mockStreamMixedBatch(server)
		return err
	}

	testID := id.NewIdFromString("test", id.Generic, t)
	gw := gateway.StartGateway(testID, servReceiverAddress, receiverImpl,
		certData, keyData, gossip.DefaultManagerFlags())

	// Init sender
	senderAddress := getNextServerAddress()
	server := StartNode(testID, senderAddress, 0,
		NewImplementation(), nil, nil)

	// Reset TLS-related global variables
	defer gw.Shutdown()
	defer server.Shutdown()

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

	mockBatch := &mixmessages.CompletedBatch{
		RoundID:     roundId,
	}

	for i := uint32(0); i < batchSize; i++ {
		mockBatch.Slots = append(mockBatch.Slots,
			&mixmessages.Slot{
				Index:    i,
				PayloadA: []byte{byte(i)},
			})
	}

	err = server.DownloadMixedBatch(host, batchInfo, mockBatch)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if receivedSlots == nil {
		t.Fatalf("Did not receive any slots")
	}

	if len(mockBatch.Slots) != len(receivedSlots) {
		t.Fatalf("Did not receive expected amount of slots." +
			"\n\tExpected: %d" +
			"\n\tReceived: %v", len(mockBatch.Slots), len(receivedSlots))
	}
}


func mockStreamMixedBatch(server mixmessages.Gateway_DownloadMixedBatchServer) ([]*mixmessages.Slot,
	error) {
	// Get header from stream
	batchInfo, err := gateway.GetMixedBatchStreamHeader(server)
	if err != nil {
		return nil, err
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

			return slots, err
		}

		// If we have another error, return err
		if err != nil {
			return slots, err
		}

		// Store slot received into temporary buffer
		slots = append(slots, slot)
	}

}
