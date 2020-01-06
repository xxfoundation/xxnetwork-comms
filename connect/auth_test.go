////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"bytes"
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
	err := c.setPrivateKey(testkeys.LoadFromPath(key))
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

	msg := &pb.AuthenticatedMessage{
		ID:        id,
		Signature: nil,
		Token:     token,
		Message:   nil,
	}

	auth := pc.AuthenticatedReceiver(msg)
	if !auth.IsAuthenticated {
		t.Errorf("Failed")
	}
}

// Happy path
func TestProtoComms_GenerateToken(t *testing.T) {
	comm := ProtoComms{
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	token, ok := comm.tokens.Load(string(tokenBytes))
	if !ok || token == nil {
		t.Errorf("Unable to find token stored in internal map")
	}
}

// Happy path
func TestProtoComms_PackAuthenticatedMessage(t *testing.T) {
	testServerId := "test12345"
	comm := ProtoComms{
		id:            testServerId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	testId := "test"
	host, err := NewHost(testId, testId, nil, false, true)
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.token = tokenBytes

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, false)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}
	if bytes.Compare(msg.Token, tokenBytes) != 0 || msg.ID != testServerId {
		t.Errorf("Expected packed message to have correct ID and Token: %+v",
			msg)
	}
}

// Happy path
func TestProtoComms_ValidateToken(t *testing.T) {
	testId := "test"
	comm := ProtoComms{
		id:            testId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	err := comm.setPrivateKey(testkeys.LoadFromPath(testkeys.GetNodeKeyPath()))
	if err != nil {
		t.Errorf("Expected to set private key: %+v", err)
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	host, err := comm.AddHost(testId, testId, pub, false, true)
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.token = tokenBytes

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, true)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}

	err = comm.ValidateToken(msg)
	if err != nil {
		t.Errorf("Expected to validate token: %+v", err)
	}
}

// Happy path
//func TestProtoComms_clientHandshake(t *testing.T) {
//	// Server Setup
//	serverAddr := "0.0.0.0:6850"
//	testId := "test"
//	comm := registration.StartRegistrationServer(serverAddr,
//		registration.NewImplementation(), nil, nil)
//	defer comm.Shutdown()
//
//	priv := testkeys.LoadFromPath(testkeys.GetNodeKeyPath())
//	err := comm.SetPrivateKey(priv)
//	if err != nil {
//		t.Errorf("Expected to set private key: %+v", err)
//	}
//	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
//
//	host, err := comm.AddHost(testId, serverAddr, pub, false, true)
//	if err != nil {
//		t.Errorf("Expected to successfully add host: %+v", err)
//	}
//
//	// Client setup
//	commSender := &client.Comms{}
//	err = commSender.SetPrivateKey(priv)
//	if err != nil {
//		t.Errorf("Expected to set private key: %+v", err)
//	}
//	hostSender, err := comm.AddHost(testId, serverAddr, pub, false, true)
//	if err != nil {
//		t.Errorf("Expected to successfully add host: %+v", err)
//	}
//	_, err = commSender.RequestNdf(hostSender, &pb.NDFHash{})
//	if err != nil {
//		t.Errorf("Expected to send successfully: %+v", err)
//	}
//}
