////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendCheckMessages
func TestSendCheckMessages(t *testing.T) {
	_, err := SendCheckMessages(SERVER_ADDRESS, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("CheckMessages: Error received: %s", err)
	}
}

// Smoke test SendGetMessage
func TestSendGetMessage(t *testing.T) {
	_, err := SendGetMessage(SERVER_ADDRESS, &pb.ClientPollMessage{})
	if err != nil {
		t.Errorf("GetMessage: Error received: %s", err)
	}
}
