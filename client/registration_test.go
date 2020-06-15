///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package client

import (
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/registration"
	"gitlab.com/elixxir/primitives/id"
	"testing"
)

// Smoke test SendRegistrationMessage
func TestSendRegistrationMessage(t *testing.T) {
	GatewayAddress := getNextAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	rg := registration.StartRegistrationServer(testId, GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	var manager connect.Manager

	host, err := manager.AddHost(testId, GatewayAddress, nil, false, false)
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
	GatewayAddress := getNextAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	rg := registration.StartRegistrationServer(testId, GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	var manager connect.Manager

	host, err := manager.AddHost(testId, GatewayAddress, nil, false, false)
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
	GatewayAddress := getNextAddress()
	testId := id.NewIdFromString("test", id.Generic, t)
	clientId := id.NewIdFromString("client", id.Generic, t)

	rg := registration.StartRegistrationServer(testId, GatewayAddress,
		registration.NewImplementation(), nil, nil)
	defer rg.Shutdown()
	c, err := NewClientComms(clientId, nil, nil, nil)
	if err != nil {
		t.Errorf("Can't create client comms: %+v", err)
	}
	var manager connect.Manager

	host, err := manager.AddHost(testId, GatewayAddress, nil, false, false)
	if err != nil {
		t.Errorf("Unable to call NewHost: %+v", err)
	}

	_, err = c.RequestNdf(host, &pb.NDFHash{})

	if err != nil {
		t.Errorf("RequestNdf: Error recieved: %s", err)
	}
}
