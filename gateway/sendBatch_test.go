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
	gateway := StartGateway(GatewayAddress, NewImplementation(), "", "")
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		"", "")
	defer gateway.Shutdown()
	defer server.Shutdown()
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID, &connect.ConnectionInfo{Address: ServerAddress})

	msgs := &pb.Batch{}
	err := gateway.PostNewBatch(connID, msgs)
	if err != nil {
		t.Errorf("PostNewBatch: Error received: %s", err)
	}
}

// Smoke Test GetBufferInfo
func TestGetRoundBufferInfo(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), "", "")
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		"", "")
	defer gateway.Shutdown()
	defer server.Shutdown()
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID, &connect.ConnectionInfo{Address: ServerAddress})

	bufSize, err := gateway.GetRoundBufferInfo(connID)
	if err != nil {
		t.Errorf("GetRoundBufferInfo: Error received: %s", err)
	}
	if bufSize != 0 {
		t.Errorf("GetRoundBufferInfo: Unexpected buffer size.")
	}
}
