////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"testing"
)

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gwShutDown := StartGateway(GatewayAddress, NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendRequestNonceMessage(ServerAddress, &pb.RequestNonceMessage{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	gwShutDown := StartGateway(GatewayAddress, NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendConfirmNonceMessage(ServerAddress, &pb.ConfirmNonceMessage{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}
