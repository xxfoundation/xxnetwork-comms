////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"testing"
)

// Smoke test PostNewBatch
func TestPostNewBatch(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	msgs := &pb.Batch{}
	err = gateway.PostNewBatch(host, msgs)
	if err != nil {
		t.Errorf("PostNewBatch: Error received: %s", err)
	}
}

// Smoke Test GetBufferInfo
func TestGetRoundBufferInfo(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	bufSize, err := gateway.GetRoundBufferInfo(&pb.Ping{}, host)
	if err != nil {
		t.Errorf("GetRoundBufferInfo: Error received: %s", err)
	}
	if bufSize.RoundBufferSize != 0 {
		t.Errorf("GetRoundBufferInfo: Unexpected buffer size.")
	}
}

// Smoke test GetCompletedBatch
func TestGetCompletedBatch(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	batch, err := gateway.GetCompletedBatch(&pb.Ping{}, host)
	if err != nil {
		t.Errorf("GetCompletedBatch: Error received: %s", err)
	}
	// The mock server doesn't have any batches ready,
	// so it should return either a nil slice of slots,
	// or a slice with no slots in it.
	if len(batch.Slots) != 0 {
		t.Errorf("GetCompletedBatch: Expected batch with no slots")
	}
}
