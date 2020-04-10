////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package mixmessages

import "gitlab.com/elixxir/primitives/current"

// GetState gets the state of the node
func (m *PermissioningPoll) GetCurrentActivityState() current.Activity {
	if m != nil {
		return current.Activity(m.Activity)
	}
	return current.NOT_STARTED
}
