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

// Smoke test RequestContactList
func TestRequestContactList(t *testing.T) {
	_, err := RequestContactList(SERVER_ADDRESS, &pb.ContactPoll{})
	if err != nil {
		t.Errorf("RequestContactList: Error received: %s", err)
	}
}
