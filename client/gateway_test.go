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
	"testing"
)

type MockID string
func (m MockID) String() string {
	return string(m)
}

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), "", "")
	defer gw.Shutdown()
	var c Client
	id := MockID("clientToGateway")
	c.ConnectToGateway(id, &connect.ConnectionInfo{
		Address:gatewayAddress,
	})

	err := c.SendPutMessage(id, &pb.Batch{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), "", "")
	connectionID := MockID("clientToGateway")
	var c Client
	c.ConnectToGateway(connectionID, &connect.ConnectionInfo{
		Address: gatewayAddress,
	})
	defer gw.Shutdown()

	_, err := c.SendCheckMessages(connectionID, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), "", "")
	connectionID := MockID("clientToGateway")
	var c Client
	c.ConnectToGateway(connectionID, &connect.ConnectionInfo{
		Address: gatewayAddress,
	})
	defer gw.Shutdown()

	_, err := c.SendGetMessage(connectionID, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), "", "")
	defer gw.Shutdown()
	connectionID := MockID("clientToGateway")
	var c Client
	c.ConnectToGateway(connectionID, &connect.ConnectionInfo{
		Address: gatewayAddress,
	})

	_, err := c.SendRequestNonceMessage(connectionID, &pb.NonceRequest{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), "", "")
	defer gw.Shutdown()
	connectionID := MockID("clientToGateway")
	var c Client
	c.ConnectToGateway(connectionID, &connect.ConnectionInfo{
		Address: gatewayAddress,
	})

	_, err := c.SendConfirmNonceMessage(connectionID, &pb.DSASignature{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}
