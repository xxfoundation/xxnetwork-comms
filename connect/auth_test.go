///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"bytes"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/elixxir/comms/testkeys"
	"gitlab.com/elixxir/crypto/csprng"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/crypto/xx"
	"gitlab.com/elixxir/primitives/id"
	token "gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
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

	message := pb.Ack{
		Error: "test",
	}

	wrappedMessage, err := ptypes.MarshalAny(&message)
	if err != nil {
		t.Errorf("Error converting to Any type: %+v", err)
	}

	signature, err := c.SignMessage(wrappedMessage)
	if err != nil {
		t.Errorf("Error signing message: %+v", err)
	}

	host := &Host{
		rsaPublicKey: pub,
	}

	err = c.VerifyMessage(wrappedMessage, signature, host)
	if err != nil {
		t.Errorf("Error verifying signature")
	}
}

func TestProtoComms_AuthenticatedReceiver(t *testing.T) {
	// Create comm object
	pc := ProtoComms{
		Manager:       Manager{},
		tokens:        token.NewMap(),
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	// Create id and token
	testID := id.NewIdFromString("testSender", id.Node, t)
	tkn := []byte("testtoken")

	// Add host
	_, err := pc.AddHost(testID, "", nil, false, true)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Get host
	h, _ := pc.GetHost(testID)
	h.receptionToken.Set(tkn)

	msg := &pb.AuthenticatedMessage{
		ID:        testID.Marshal(),
		Signature: nil,
		Token:     tkn,
		Message:   nil,
	}

	// Try the authenticated received
	auth, err := pc.AuthenticatedReceiver(msg)
	if err != nil {
		t.Errorf("AuthenticatedReceiver() produced an error: %v", err)
	}

	if !auth.IsAuthenticated {
		t.Errorf("Failed: authenticated receiver")
	}

	// Compare the tokens
	if !bytes.Equal(auth.Sender.receptionToken.Get(), tkn) {
		t.Errorf("Tokens do not match! \n\tExptected: " +
			"%+v\n\tReceived: %+v", tkn, auth.Sender.receptionToken)
	}
}

// Error path
func TestProtoComms_AuthenticatedReceiver_BadId(t *testing.T) {
	// Create comm object
	pc := ProtoComms{
		Manager:       Manager{},
		tokens:        token.NewMap(),
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	// Create id and token
	testID := id.NewIdFromString("testSender", id.Node, t)
	tkn := []byte("testtoken")

	// Add host
	_, err := pc.AddHost(testID, "", nil, false, true)
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Get host
	h, _ := pc.GetHost(testID)
	h.receptionToken.Set(tkn)

	badId := []byte("badID")

	msg := &pb.AuthenticatedMessage{
		ID:        badId,
		Signature: nil,
		Token:     tkn,
		Message:   nil,
	}

	// Try the authenticated received
	_, err = pc.AuthenticatedReceiver(msg)
	if err != nil {
		return
	}

	t.Errorf("Expected error path!"+
		"Should not be able to marshal a message with id: %v", badId)

}

// Happy path
func TestProtoComms_GenerateToken(t *testing.T) {
	comm := ProtoComms{
		tokens:token.NewMap(),
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
	}
	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	tkn, err := token.Unmarshal(tokenBytes)
	if err!=nil{
		t.Errorf("Should be able to unmarshal token: %s", err)
	}

	ok := comm.tokens.Validate(tkn)
	if !ok{
		t.Errorf("Unable to validate token")
	}
}

// Happy path
func TestProtoComms_PackAuthenticatedMessage(t *testing.T) {
	testServerId := id.NewIdFromString("test12345", id.Node, t)
	comm := ProtoComms{
		Id:            testServerId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
		tokens:token.NewMap(),
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	testId := id.NewIdFromString("test", id.Node, t)

	host, err := NewHost(testId, "test", nil, false, true)
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tokenBytes)

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, false)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}
	// Compare the tokens and id's
	if bytes.Compare(msg.Token, tokenBytes) != 0 || !bytes.Equal(msg.ID, testServerId.Marshal()) {
		t.Errorf("Expected packed message to have correct ID and Live: %+v",
			msg)
	}
}

// Happy path
func TestProtoComms_ValidateToken(t *testing.T) {
	testId := id.NewIdFromString("test", id.Node, t)
	comm := ProtoComms{
		Id:            testId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
		tokens:token.NewMap(),
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
	host, err := comm.AddHost(testId, "test", pub, false, true)
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tokenBytes)

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, true)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}

	msg.Client.PublicKey = string(pub)

	// Check the token
	err = comm.ValidateToken(msg)
	if err != nil {
		t.Errorf("Expected to validate token: %+v", err)
	}

	if !bytes.Equal(msg.Token, host.transmissionToken.Get()) {
		t.Errorf("Message token doesn't match message's token! "+
			"Expected: %+v"+
			"\n\tReceived: %+v", host.transmissionToken, msg.Token)
	}
}

// Error Path
func TestProtoComms_ValidateToken_BadId(t *testing.T) {
	testId := id.NewIdFromString("test", id.Node, t)
	comm := ProtoComms{
		Id:            testId,
		LocalServer:   nil,
		ListeningAddr: "",
		privateKey:    nil,
		tokens:token.NewMap(),
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
	host, err := comm.AddHost(testId, "test", pub, false, true)
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tokenBytes)

	tokenMsg := &pb.AssignToken{
		Token: tokenBytes,
	}

	msg, err := comm.PackAuthenticatedMessage(tokenMsg, host, true)
	if err != nil {
		t.Errorf("Expected no error packing authenticated message: %+v", err)
	}

	// Assign message a bad id
	badId := []byte("badID")
	msg.ID = badId

	msg.Client.PublicKey = string(pub)

	// Check the token
	err = comm.ValidateToken(msg)
	if err != nil {
		return
	}

	t.Errorf("Expected error path!"+
		"Should not be able to marshal a message with id: %v", badId)

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
	uid, err := xx.NewID(pubKey, salt, id.User)
	if err != nil {
		t.Errorf("Could not generate user ID: %+v", err)
	}
	testId := uid.String()
	// ------

	// Now we set up the client comms object
	comm := ProtoComms{
		Id:            uid,
		ListeningAddr: "",
		tokens:token.NewMap(),
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
	host, err := newDynamicHost(uid, rsa.CreatePublicKeyPem(pubKey))
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tokenBytes)
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
	host, ok := comm.GetHost(uid)
	if !ok {
		t.Errorf("Expected dynamic auth to add host %s!", testId)
	}
	if !host.IsDynamicHost() {
		t.Errorf("Expected host to be dynamic!")
	}

	if !bytes.Equal(msg.Token, host.receptionToken.Get()) {
		t.Errorf("Message token doesn't match message's token! "+
			"Expected: %+v"+
			"\n\tReceived: %+v", host.receptionToken, msg.Token)
	}

}

func TestProtoComms_DisableAuth(t *testing.T) {
	testId := id.NewIdFromString("test", id.Node, t)
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
