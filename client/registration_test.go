////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"testing"
)

// Smoke test SendRegistrationMessage
func TestSendRegistrationMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	rg := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, GatewayAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendRegistrationMessage(host, &pb.UserRegistration{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Smoke test SendCheckClientVersion
func TestSendCheckClientVersionMessage(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	rg := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, GatewayAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.SendGetCurrentClientVersionMessage(host)
	if err != nil {
		t.Errorf("CheckClientVersion: Error received: %s", err)
	}
}

//Smoke test RequestNdf
func TestSendGetUpdatedNDF(t *testing.T) {
	GatewayAddress := getNextGatewayAddress()
	rg := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	var c Comms
	var manager connect.Manager

	testId := "test"
	host, err := manager.AddHost(testId, GatewayAddress, nil, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestNdf(host, &pb.NDFHash{})

	if err != nil {
		t.Errorf("RequestNdf: Error recieved: %s", err)
	}
}
