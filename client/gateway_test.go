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
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms

	err := c.SendPutMessage(&connect.Host{
		address:        gatewayAddress,
		certificate:    nil,
		disableTimeout: false,
	}, &pb.Slot{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()

	_, err := c.SendCheckMessages(&connect.Host{
		address:        gatewayAddress,
		certificate:    nil,
		disableTimeout: false,
	}, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()

	_, err := c.SendGetMessage(&connect.Host{
		address:        gatewayAddress,
		certificate:    nil,
		disableTimeout: false,
	}, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms

	_, err := c.SendRequestNonceMessage(&connect.Host{
		address:        gatewayAddress,
		certificate:    nil,
		disableTimeout: false,
	}, &pb.NonceRequest{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	gatewayAddress := getNextGatewayAddress()
	gw := gateway.StartGateway(gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms

	_, err := c.SendConfirmNonceMessage(&connect.Host{
		address:        gatewayAddress,
		certificate:    nil,
		disableTimeout: false,
	},
		&pb.RequestRegistrationConfirmation{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %s", err)
	}
}
