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
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendPutMessage(host, &pb.GatewaySlot{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	var c Comms
	defer gw.Shutdown()
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
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

// Smoke test SendRequestNonceMessage
func TestSendRequestNonceMessage(t *testing.T) {
	gatewayAddress := getNextAddress()
	testID := id.NewIdFromString("test", id.Gateway, t)
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
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
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
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
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
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
	gw := gateway.StartGateway(testID, gatewayAddress,
		gateway.NewImplementation(), nil, nil, gossip.DefaultManagerFlags())
	defer gw.Shutdown()
	var c Comms
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestMessages(host,
		&pb.MessageRequest{})
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
	var c Comms
	var manager connect.Manager

	host, err := manager.AddHost(testID, gatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestBloom(host,
		&pb.GetBloom{})
	if err != nil {
		t.Errorf("SendPoll: Error received: %+v", err)
	}
}
