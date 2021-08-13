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
	"git.xx.network/xx_network/comms/messages"
	"git.xx.network/xx_network/primitives/id"
	"testing"
)

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()
	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(testID, ServerAddress, 0, node.NewImplementation(),
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

	RSASignature := &messages.RSASignature{
		Signature: []byte("test"),
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
	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(testID, ServerAddress, 0, node.NewImplementation(),
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

	reg := &pb.RequestRegistrationConfirmation{UserID: testID.Bytes()}
	reg.NonceSignedByClient = &messages.RSASignature{}
	_, err = gateway.SendConfirmNonceMessage(host, reg)
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}

func TestPoll(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	ServerAddress := getNextServerAddress()

	testID := id.NewIdFromString("test", id.Generic, t)
	gateway := StartGateway(testID, GatewayAddress, NewImplementation(), nil,
		nil, gossip.DefaultManagerFlags())
	server := node.StartNode(testID, ServerAddress, 0, node.NewImplementation(),
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

	_, err = gateway.SendPoll(host, &pb.ServerPoll{})
	if err != nil {
		t.Errorf("TestDemandNdf: Error received: %s", err)
	}
}
