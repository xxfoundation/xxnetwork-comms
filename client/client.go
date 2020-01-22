////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level Client comms object

package client

import (
	"gitlab.com/elixxir/comms/connect"
)

// Client object used to implement endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms
}

// Returns a Comms object with given attributes
func NewClientComms(id string) *Comms {
	return &Comms{connect.ProtoComms{
		Id: id,
	}}
}
