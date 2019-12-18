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
	"sync"
	"testing"
)

func TestSignVerify(t *testing.T) {

	c := *new(ProtoComms)

	key := testkeys.GetNodeKeyPath()
	err := c.SetPrivateKey(testkeys.LoadFromPath(key))
	if err != nil {
		t.Errorf("Error setting private key: %+v", err)
	}

	private := c.GetPrivateKey()
	pub := private.Public().(*rsa.PublicKey)

	message := pb.NDF{
		Ndf: []byte("test"),
	}

	wrappedMessage, err := ptypes.MarshalAny(&message)
	if err != nil {
		t.Errorf("Error converting to Any type: %+v", err)
	}

	signature, err := c.signMessage(wrappedMessage)
	if err != nil {
		t.Errorf("Error signing message: %+v", err)
	}

	host := &Host{
		rsaPublicKey: pub,
	}

	err = c.verifyMessage(&pb.AuthenticatedMessage{
		ID:        "",
		Signature: signature,
		Token:     nil,
		Message:   wrappedMessage,
	}, host)
	if err != nil {
		t.Errorf("Error verifying signature")
	}
}

func TestProtoComms_AuthenticatedReceiver(t *testing.T) {
	pc := ProtoComms{
		Manager:       Manager{},
		tokens:        sync.Map{},
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	id := "testsender"
	token := []byte("testtoken")

	_, err := pc.AddHost(id, "", nil, false, true)
	if err != nil {
		t.Errorf("uh oh")
	}
	h, _ := pc.GetHost(id)
	h.token = token

	var authenticatedTokens sync.Map
	msg := pb.AuthenticatedMessage{
		ID:        id,
		Signature: nil,
		Token:     token,
		Message:   nil,
	}

	auth := pc.AuthenticatedReceiver(msg, authenticatedTokens)
	if !auth.IsAuthenticated {
		t.Errorf("Failed")
	}
}
