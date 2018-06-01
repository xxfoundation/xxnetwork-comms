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

// Smoke test SendClientPoll
func TestSendClientPoll(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress, gateway.NewImplementation())
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation())
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendClientPoll(ServerAddress, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("SendClientPoll: Error received: %s", err)
	}
}

// Smoke test SendRegistrationPoll
func TestSendRegistrationPoll(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress, gateway.NewImplementation())
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation())
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendRegistrationPoll(ServerAddress, &pb.RegistrationPoll{})
	if err != nil {
		t.Errorf("SendRegistrationPoll: Error received: %s", err)
	}
}

// Smoke test SendMessagetoSender
func TestSendMessageToServer(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress, gateway.NewImplementation())
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation())
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendMessageToServer(ServerAddress, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("SendMessageToServer: Error received: %s", err)
	}
}
