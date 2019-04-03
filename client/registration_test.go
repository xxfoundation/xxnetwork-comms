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
	rgShutDown := registration.StartRegistrationServer(GatewayAddress,
		registration.NewImplementation(), "", "")
	defer rgShutDown()

	_, err := SendRegistrationMessage(GatewayAddress, "", "",
		&pb.RegisterUserMessage{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}
