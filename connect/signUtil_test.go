////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"github.com/golang/protobuf/ptypes"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"testing"
)

func TestSignVerify(t *testing.T) {

	c := *new(ConnectionManager)

	err := c.SetPrivateKey(testkeys.GetNodeKeyPath())
	if err != nil {
		t.Errorf("Error setting private key: %+v", err)
	}

	private := c.GetPrivateKey()
	pub := private.Public().(*rsa.PublicKey)
	c.SetPublicKey(pub)

	message := pb.NodeTopology{
		Topology: []*pb.NodeInfo{
			{
				Id:        []byte("test"),
				Index:     uint32(3),
				IpAddress: "0.0.0.0",
			},
		},
	}

	wrappedMessage, err := ptypes.MarshalAny(&message)
	if err != nil {
		t.Errorf("Error converting to Any type: %+v", err)
	}

	signed, err := c.SignMessage(wrappedMessage)
	if err != nil {
		t.Errorf("Error signing message: %+v", err)
	}

	verified := pb.NodeTopology{}
	err = c.VerifySignature(signed, &verified)
	if err != nil {
		t.Errorf("Error verifying signature")
	}

	t.Logf("%+v", verified)
}
