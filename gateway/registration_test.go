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
	gateway := StartGateway(GatewayAddress, NewImplementation(), nil, nil)
	server := node.StartNode(ServerAddress, node.NewImplementation(),
		nil, nil)
	defer gateway.Shutdown()
	defer server.Shutdown()

	RSASignature := &pb.RSASignature{
		Signature: []byte{},
	}

	_, err := gateway.SendRequestNonceMessage(&connect.Host{
		address:        ServerAddress,
		certificate:    nil,
		disableTimeout: false,
	},
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

	reg := &pb.RequestRegistrationConfirmation{}
	reg.NonceSignedByClient = &pb.RSASignature{}
	_, err := gateway.SendConfirmNonceMessage(&connect.Host{
		address:        ServerAddress,
		certificate:    nil,
		disableTimeout: false,
	}, reg)
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

	_, err := gateway.PollSignedCerts(&connect.Host{
		address:        ServerAddress,
		certificate:    nil,
		disableTimeout: false,
	}, &pb.Ping{})
	if err != nil {
		t.Errorf("SendGetSignedCertMessage: Error received: %s", err)
	}
}
