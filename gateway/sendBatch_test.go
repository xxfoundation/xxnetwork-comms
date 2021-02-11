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
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Smoke test PostNewBatch
func TestPostNewBatch(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(testID, ServerAddress, 0, node.NewImplementation(),
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
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(nodeID, ServerAddress, 0, node.NewImplementation(),
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
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(nodeID, ServerAddress, 0, node.NewImplementation(),
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

// Smoke test SendShareMessages
func TestComms_SendShareMessages(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	GatewayAddress2 := getNextGatewayAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	testID2 := id.NewIdFromString("test2", id.Gateway, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	gateway2 := StartGateway(testID2, GatewayAddress2, NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gateway.Shutdown()
	defer gateway2.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, GatewayAddress2, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	err = gateway.SendShareMessages(host, &pb.RoundMessages{})
	if err != nil {
		t.Errorf("ShareMessages: Error received: %s", err)
	}
}
