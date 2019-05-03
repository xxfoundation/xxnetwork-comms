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
		registration.NewImplementation(), "", "")
	defer rg.Shutdown()
	connID := MockID("clientToRegistration")
	var c ClientComms
	c.ConnectToRegistration(connID, &connect.ConnectionInfo{Address: GatewayAddress})

	_, err := c.SendRegistrationMessage(connID, &pb.UserRegistration{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}
