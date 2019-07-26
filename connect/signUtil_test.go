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

	key := testkeys.GetNodeKeyPath()
	err := c.SetPrivateKey(testkeys.LoadFromPath(key))
	if err != nil {
		t.Errorf("Error setting private key: %+v", err)
	}

	private := c.GetPrivateKey()
	pub := private.Public().(*rsa.PublicKey)

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

	signed, err := c.SignMessage(wrappedMessage, "test_id")
	if err != nil {
		t.Errorf("Error signing message: %+v", err)
	}

	verified := pb.NodeTopology{}
	err = ptypes.UnmarshalAny(signed.Message, &verified)
	if err != nil {
		t.Errorf("Failed to unmarshal generic message, check your input message type: %+v", err)
	}

	err = c.VerifySignature(signed, &verified, pub)
	if err != nil {
		t.Errorf("Error verifying signature")
	}

	if len(verified.Topology) != 1 && string(verified.Topology[0].Id) != "test" {
		t.Errorf("Message contents do not match original: %+v", verified)
	}
}
