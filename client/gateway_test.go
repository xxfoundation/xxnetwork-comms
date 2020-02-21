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

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	err = c.SendPutMessage(host, &pb.Slot{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendCheckMessages(host, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()

	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendGetMessage(host, &pb.ClientRequest{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendRequestNonceMessage(host, &pb.NonceRequest{})
	if err != nil {
		t.Errorf("SendRequestNonceMessage: Error received: %s", err)
	}
}

// Smoke test SendConfirmNonceMessage
func TestSendConfirmNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendConfirmNonceMessage(host,
		&pb.RequestRegistrationConfirmation{})
	if err != nil {
		t.Errorf("SendConfirmNonceMessage: Error received: %+v", err)
	}
}

// Smoke test SendPoll
func TestComms_SendPoll(t *testing.T) {
	gatewayAddress := getNextAddress()
	gw := gateway.StartGateway("test", gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendPoll(host,
		&pb.GatewayPoll{
			Partial: &pb.NDFHash{
				Hash: make([]byte, 0),
			},
			LastRealtimeKnown: 0,
			LastKilledRound:   0,
			LastMessageID:     "",
		})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}
