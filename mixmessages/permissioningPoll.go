////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
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
