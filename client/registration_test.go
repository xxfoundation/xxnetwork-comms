////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
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
	connID := MockID("clientToRegistration")
	var c ClientComms
	c.ConnectToRemote(connID, GatewayAddress, nil, false)

	_, err := c.SendRegistrationMessage(connID, &pb.UserRegistration{})
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
	connID := MockID("clientToRegistration")
	var c ClientComms
	c.ConnectToRemote(connID, GatewayAddress, nil, false)

	_, err := c.SendCheckClientVersionMessage(connID, &pb.ClientVersion{})
	if err != nil {
		t.Errorf("CheckClientVersion: Error received: %s", err)
	}
}
