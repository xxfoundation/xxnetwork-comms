////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package client

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendClientPoll
func TestSendClientPoll(t *testing.T) {
	_, err := SendClientPoll(ServerAddress, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("SendClientPoll: Error received: %s", err)
	}
}

// Smoke test SendRegistrationPoll
func TestSendRegistrationPoll(t *testing.T) {
	_, err := SendRegistrationPoll(ServerAddress, &pb.RegistrationPoll{})
	if err != nil {
		t.Errorf("SendRegistrationPoll: Error received: %s", err)
	}
}

// Smoke test SendMessagetoSender
func TestSendMessageToServer(t *testing.T) {
	_, err := SendMessageToServer(ServerAddress, &pb.CmixMessage{})
	if err != nil {
		t.Errorf("SendMessageToServer: Error received: %s", err)
	}
}
