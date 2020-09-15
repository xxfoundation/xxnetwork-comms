///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"testing"
)

// Smoke test PostNewBatch
func TestPostNewBatch(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode(testID, ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	nodeID := id.NewIdFromString("test", id.Node, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode(nodeID, ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	bufSize, err := gateway.GetRoundBufferInfo(host)
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	nodeID := id.NewIdFromString("test", id.Node, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode(nodeID, ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, ServerAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	batch, err := gateway.GetCompletedBatch(host)
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
