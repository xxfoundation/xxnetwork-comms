////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/privategrity/comms/gateway"
	pb "gitlab.com/privategrity/comms/mixmessages"
	"gitlab.com/privategrity/comms/node"
	"testing"
)

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	err := SendPutMessage(GatewayAddress, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendCheckMessages(GatewayAddress, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendGetMessage(GatewayAddress, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}
