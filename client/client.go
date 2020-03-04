////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level Client comms object

package client

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
)

// Client object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
}

// Returns a Comms object with given attributes
func NewClientComms(id string, pubKeyPem, privKeyPem, salt []byte) (*Comms, error) {
	pc, err := connect.CreateCommClient(id, pubKeyPem, privKeyPem, salt)
	if err != nil {
		return nil, errors.Errorf("Unable to create Client comms: %+v", err)
	}
	return &Comms{pc}, nil
}
