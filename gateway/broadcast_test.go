///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "git.xx.network/elixxir/comms/mixmessages"
	"git.xx.network/elixxir/comms/node"
	"git.xx.network/xx_network/comms/connect"
	"git.xx.network/xx_network/comms/gossip"
	"git.xx.network/xx_network/primitives/id"
	"testing"
)

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
