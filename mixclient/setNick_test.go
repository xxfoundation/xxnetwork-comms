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

// Smoke test SetNick
func TestSetNick(t *testing.T) {
	_, err := SetNick(SERVER_ADDRESS, &pb.Contact{})
	if err != nil {
		t.Errorf("SetNick: Error received: %s", err)
	}
}
