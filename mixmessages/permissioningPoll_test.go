////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import "testing"

func TestPermissioningPoll_GetActivity(t *testing.T) {
	expected := uint32(45)
	testRoundInfo := &PermissioningPoll{
		Activity: expected,
	}

	received := testRoundInfo.GetActivity()

	if received != expected {
		t.Errorf("Received does not match expected for getter function! "+
			"Expected: %+v \n\t"+
			"Received: %+v", expected, received)
	}
}
