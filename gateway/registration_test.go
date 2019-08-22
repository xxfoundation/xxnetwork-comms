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
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID, ServerAddress, nil, true)

	RSASignature := &pb.RSASignature{
		Signature: []byte{},
	}

	_, err := gateway.SendRequestNonceMessage(connID,
		&pb.NonceRequest{ClientSignedByServer: RSASignature,
			RequestSignature: RSASignature})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID, ServerAddress, nil, true)

	reg := &pb.RequestRegistrationConfirmation{}
	reg.NonceSignedByClient = &pb.RSASignature{}
	_, err := gateway.SendConfirmNonceMessage(connID, reg)
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}

func TestPollSignedCerts(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()

	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(), nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	connID := MockID("gatewayToServer")
	gateway.ConnectToNode(connID, ServerAddress, nil, true)

	_, err := gateway.PollSignedCerts(connID, &pb.Ping{})
	if err != nil {
		t.Errorf("SendGetSignedCertMessage: Error received: %s", err)
	}
}
