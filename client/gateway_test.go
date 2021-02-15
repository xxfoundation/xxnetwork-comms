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
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/gossip"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

// Smoke test SendGetMessage
func TestSendPutMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms
	manager := connect.NewManagerTesting(t)

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false
	host, err := manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendPutMessage(host, &pb.GatewaySlot{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
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
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
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
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
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
		})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}

// Smoke test RequestMessages
func TestComms_RequestMessages(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	pk := testkeys.LoadFromPath(testkeys.GetGatewayKeyPath())

	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	c, err := NewClientComms(testID, nil, pk, nil)
	if err != nil {
		t.Errorf("Could not start client: %v", err)
	}
	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false

	host, err := c.Manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestMessages(host,
		&pb.GetMessages{})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}

// Smoke test RequestHistoricalRounds
func TestComms_RequestHistoricalRounds(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	pk := testkeys.LoadFromPath(testkeys.GetGatewayKeyPath())

	c, err := NewClientComms(testID, nil, pk, nil)
	if err != nil {
		t.Errorf("Could not start client: %v", err)
	}

	params := connect.GetDefaultHostParams()
	params.AuthEnabled = false

	host, err := c.Manager.AddHost(testID, gatewayAddress, nil, params)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestHistoricalRounds(host,
		&pb.HistoricalRounds{})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}
