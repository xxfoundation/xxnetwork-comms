////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"testing"
)

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	err := SendPutMessage(GatewayAddress, "", "", &pb.Batch{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Fail test SendGetMessage
func TestSendPutMessage_Failure(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	err := SendPutMessage(ServerAddress, &pb.CmixMessage{})
	if err == nil {
		t.Errorf("PutMessage: Expected error!")
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendCheckMessages(GatewayAddress, "", "", &pb.ClientRequest{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Fail test SendCheckMessages
func TestSendCheckMessages_Failure(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendCheckMessages(ServerAddress, &pb.ClientPollMessage{})
	if err == nil {
		t.Errorf("CheckMessages: Expected error!")
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendGetMessage(GatewayAddress, "", "", &pb.ClientRequest{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}

// Fail test SendGetMessage
func TestSendGetMessage_Failure(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendGetMessage(ServerAddress, &pb.ClientPollMessage{})
	if err == nil {
		t.Errorf("GetMessage: Expected error!")
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendRequestNonceMessage(GatewayAddress, "", "",
		&pb.NonceRequest{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Fail test SendRequestNonceMessage
func TestSendRequestNonceMessage_Failure(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendRequestNonceMessage(ServerAddress, &pb.RequestNonceMessage{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Expected error!")
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendConfirmNonceMessage(GatewayAddress, "", "",
		&pb.DSASignature{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}

// Fail test SendConfirmNonceMessage
func TestSendConfirmNonceMessage_Failure(t *testing.T) {
	gwShutDown := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gwShutDown()
	defer nodeShutDown()

	_, err := SendConfirmNonceMessage(ServerAddress, &pb.ConfirmNonceMessage{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Expected error!")
	}
}
