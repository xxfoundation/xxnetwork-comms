////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package client

import (
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/node"
	"gitlab.com/elixxir/comms/registration"
	"testing"
)

// Smoke test SendRegistrationMessage
func TestSendRegistrationMessage(t *testing.T) {
	rgShutDown := registration.StartRegistrationServer(RegistrationAddress,
		registration.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer rgShutDown()
	defer nodeShutDown()

	_, err := SendRegistrationMessage(RegistrationAddress, &pb.RegisterUserMessage{})
	if err != nil {
		t.Errorf("RegistrationMessage: Error received: %s", err)
	}
}

// Fail test SendRegistrationMessage
func TestSendRegistrationMessage_Failure(t *testing.T) {
	rgShutDown := registration.StartRegistrationServer(RegistrationAddress,
		registration.NewImplementation(), "", "")
	nodeShutDown := node.StartServer(ServerAddress, node.NewImplementation(),
		"", "")
	defer rgShutDown()
	defer nodeShutDown()

	_, err := SendRegistrationMessage(ServerAddress, &pb.RegisterUserMessage{})
	if err == nil {
		t.Errorf("RegistrationMessage: Expected error!")
	}
}
