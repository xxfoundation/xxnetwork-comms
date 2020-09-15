///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/xx_network/comms/connect"
	"testing"
)

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	testUserID := id.NewIdFromString("Test User", id.User, t)
	testUserEntityID := &pb.ClientRequest{UserID: testUserID.Bytes()}

	_, err = c.SendCheckMessages(host, testUserEntityID)
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	var c Comms
	defer gw.Shutdown()

	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	testUserID := id.NewIdFromString("Test User", id.User, t)
	testUserEntityID := &pb.ClientRequest{UserID: testUserID.Bytes()}

	_, err = c.SendGetMessage(host, testUserEntityID)
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
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
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil)
	defer gw.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendPoll(host,
		&pb.GatewayPoll{
			Partial: &pb.NDFHash{
				Hash: make([]byte, 0),
			},
			LastMessageID: "",
		})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}
