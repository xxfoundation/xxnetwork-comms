////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level Client comms object

package client

import (
	"github.com/pkg/errors"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
)

// Client object used to implement endpoints and top-level comms functionality
type Comms struct {
	*connect.ProtoComms
}

// Returns a Comms object with given attributes
func NewClientComms(id *id.ID, pubKeyPem, privKeyPem, salt []byte) (*Comms, error) {
	pc, err := connect.CreateCommClient(id, pubKeyPem, privKeyPem, salt)
	if err != nil {
		return nil, errors.Errorf("Unable to create Client comms: %+v", err)
	}
	return &Comms{pc}, nil
}
