package client

import (
	"github.com/pkg/errors"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/primitives/id"
)

// Comms is an object used for top-level remote sync client calls.
type Comms struct {
	*connect.ProtoComms
}

// NewClientComms returns a Comms object with given attributes.
func NewClientComms(id *id.ID, pubKeyPem, privKeyPem, salt []byte) (*Comms, error) {
	pc, err := connect.CreateCommClient(id, pubKeyPem, privKeyPem, salt)
	if err != nil {
		return nil, errors.Errorf("Unable to create Client comms: %+v", err)
	}
	return &Comms{pc}, nil
}
