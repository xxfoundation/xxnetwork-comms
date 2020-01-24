////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level Client comms object

package client

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// Client object used to implement endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms

	// Used to store the public key used for generating Client Id
	publicKey []byte
	// Used to store the salt used for generating Client Id
	salt []byte
}

// Returns a Comms object with given attributes
func NewClientComms(id string, pubKey, salt []byte) *Comms {
	return &Comms{
		connect.ProtoComms{
			Id: id,
		},
		pubKey, salt,
	}
}

// Wrapper for PackAuthenticatedMessage that adds special client info
// to the newly-generated authenticated message
func (c *Comms) PackAuthenticatedClientMessage(msg proto.Message,
	host *connect.Host, enableSignature bool) (*pb.AuthenticatedMessage, error) {

	authMsg, err := c.PackAuthenticatedMessage(msg, host, enableSignature)
	if err != nil {
		return nil, err
	}

	authMsg.Client = &pb.ClientID{
		Salt:      c.salt,
		PublicKey: string(c.publicKey),
	}

	return authMsg, nil
}
