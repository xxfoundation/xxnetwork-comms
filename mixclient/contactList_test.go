////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package mixclient

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"testing"
)

// Smoke test SendClientGetContactList
func TestSendClientGetContactList(t *testing.T) {
	_, err := SendClientGetContactList(SERVER_ADDRESS, &pb.ContactPoll{})
	if err != nil {
		t.Errorf("SendClientGetContactList: Error received: %s", err)
	}
}
