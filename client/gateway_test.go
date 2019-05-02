////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"testing"
)

type MockID string
func (m MockID) String() string {
	return string(m)
}

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gw := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	server := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer gw.Shutdown()
	defer server.Shutdown()

	err := server.SendPutMessage(GatewayAddress, "", "", &pb.Batch{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	clientAddress := getNextServerAddress()
	gw := gateway.StartGateway(GatewayAddress,
		gateway.NewImplementation(), "", "")
	c := StartClient(clientAddress, newImplementation()
		"", "")
	connectionID := MockID("connection76")
	c.ConnectToGateway(connectionID, &connect.ConnectionInfo{
		Address: GatewayAddress,
	})
	defer gw.Shutdown()
	defer c.Shutdown()

	_, err := c.SendCheckMessages(connectionID, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
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
