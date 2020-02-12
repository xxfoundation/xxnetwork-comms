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

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	gateway := StartGateway("test", GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode("test", ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	RSASignature := &pb.RSASignature{
		Signature: []byte{},
	}

	_, err = gateway.SendRequestNonceMessage(host,
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
	gateway := StartGateway("test", GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode("test", ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	reg := &pb.RequestRegistrationConfirmation{}
	reg.NonceSignedByClient = &pb.RSASignature{}
	_, err = gateway.SendConfirmNonceMessage(host, reg)
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}

func TestPollNdf(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()

	gateway := StartGateway("test", GatewayAddress, NewImplementation(), nil,
		nil)
	server := node.StartNode("test", ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, ServerAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = gateway.DemandNdf(host)
	if err != nil {
		t.Errorf("SendGetSignedCertMessage: Error received: %s", err)
	}
}
