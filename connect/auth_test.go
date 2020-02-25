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
	"gitlab.com/elixxir/crypto/csprng"
	"gitlab.com/elixxir/crypto/registration"
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
	// Create comm object
	pc := ProtoComms{
		Manager:       Manager{},
		tokens:        sync.Map{},
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	// Create id and token
	id := "testsender"
	token := []byte("testtoken")

	// Add host
	_, err := pc.AddHost(id, "", nil, false, true)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Get host
	h, _ := pc.GetHost(id)
	h.receptionToken = token

	msg := &pb.AuthenticatedMessage{
		ID:        id,
		Signature: nil,
		Token:     token,
		Message:   nil,
	}

	// Try the authenticated received
	auth := pc.AuthenticatedReceiver(msg)
	if !auth.IsAuthenticated {
		t.Errorf("Failed: authenticated receiver")
	}

	// Compare the tokens
	if !bytes.Equal(auth.Sender.receptionToken, token) {
		t.Errorf("Tokens do not match! \n\tExptected: %+v\n\tReceived: %+v", token, auth.Sender.receptionToken)
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
		Id:            testServerId,
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
	host.transmissionToken = tokenBytes

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, false)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}
	// Compare the tokens and id's
	if bytes.Compare(msg.Token, tokenBytes) != 0 || msg.ID != testServerId {
		t.Errorf("Expected packed message to have correct ID and Token: %+v",
			msg)
	}
}

// Happy path
func TestProtoComms_ValidateToken(t *testing.T) {
	testId := "test"
	comm := ProtoComms{
		Id:            testId,
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
	host.transmissionToken = tokenBytes

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, true)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}

	// Check the token
	err = comm.ValidateToken(msg)
	if err != nil {
		t.Errorf("Expected to validate token: %+v", err)
	}
}

// Dynamic authentication happy path (e.g. host not pre-added)
func TestProtoComms_ValidateTokenDynamic(t *testing.T) {
	// All of this is setup for UID ----
	privKey, err := rsa.GenerateKey(csprng.NewSystemRNG(), rsa.DefaultRSABitLen)
	if err != nil {
		t.Errorf("Could not generate private key: %+v", err)
	}
	pubKey := privKey.GetPublic()

	salt := []byte("0123456789ABCDEF0123456789ABCDEF")
	uid := registration.GenUserID(pubKey, salt)
	testId := uid.String()
	// ------

	// Now we set up the client comms object
	comm := ProtoComms{
		Id:            testId,
		ListeningAddr: "",
	}
	err = comm.setPrivateKey(rsa.CreatePrivateKeyPem(privKey))
	if err != nil {
		t.Errorf("Expected to set private key: %+v", err)
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	// For this test we won't addHost to Manager, we'll just create a host
	// so we can compare to the dynamic one later
	host, err := newDynamicHost(testId, rsa.CreatePublicKeyPem(pubKey))
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken = tokenBytes
	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	// Set up auth msg
	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, true)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}
	msg.Client = &pb.ClientID{
		Salt:      salt,
		PublicKey: string(rsa.CreatePublicKeyPem(pubKey)),
	}

	// Here's the method we're testing
	err = comm.ValidateToken(msg)
	if err != nil {
		t.Errorf("Expected to validate token: %+v", err)
	}

	// Check the output values behaved as expected
	host, ok := comm.GetHost(testId)
	if !ok {
		t.Errorf("Expected dynamic auth to add host %s!", testId)
	}
	if !host.IsDynamicHost() {
		t.Errorf("Expected host to be dynamic!")
	}

}

func TestProtoComms_DisableAuth(t *testing.T) {
	testId := "test"
	comm := ProtoComms{
		Id:            testId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}

	comm.DisableAuth()

	if !comm.disableAuth {
		t.Error("Auth was not disabled when DisableAuth was called")
	}
}
