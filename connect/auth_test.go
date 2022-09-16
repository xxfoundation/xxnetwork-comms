////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package connect

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/ptypes"
	token "gitlab.com/xx_network/comms/connect/token"
	pb "gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/comms/testkeys"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc/peer"
	"net"
	"testing"
	"time"
)

func TestSignVerify(t *testing.T) {

	c := *new(ProtoComms)

	key := testkeys.GetNodeKeyPath()
	err := c.setPrivateKey(testkeys.LoadFromPath(key))
	if err != nil {
		t.Errorf("Error setting private key: %+v", err)
	}

	testId := id.NewIdFromBytes([]byte("Kirby"), t)
	c.networkId = testId

	private := c.GetPrivateKey()
	pub := private.Public().(*rsa.PublicKey)

	message := pb.Ack{
		Error: "test",
	}

	wrappedMessage, err := ptypes.MarshalAny(&message)
	if err != nil {
		t.Errorf("Error converting to Any type: %+v", err)
	}

	signature, err := c.signMessage(wrappedMessage, testId)
	if err != nil {
		t.Errorf("Error signing message: %+v", err)
	}

	host := &Host{
		id:           testId,
		rsaPublicKey: pub,
	}

	err = c.verifyMessage(wrappedMessage, signature, host)
	if err != nil {
		t.Errorf("Error verifying signature")
	}
}

func TestProtoComms_AuthenticatedReceiver(t *testing.T) {
	// Create comm object
	pc := ProtoComms{
		Manager:    newManager(),
		tokens:     token.NewMap(),
		grpcServer: nil,
		privateKey: nil,
	}
	// Create id and token
	testID := id.NewIdFromString("testSender", id.Node, t)
	expectedVal := []byte("testToken")
	tkn := token.Token{}
	copy(tkn[:], expectedVal)

	// Add host
	_, err := pc.AddHost(testID, "", nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Failed to add host: %+v", err)
	}

	// Get host
	h, _ := pc.GetHost(testID)
	h.receptionToken.Set(tkn)

	msg := &pb.AuthenticatedMessage{
		ID:        testID.Marshal(),
		Signature: nil,
		Token:     tkn.Marshal(),
		Message:   nil,
	}

	// Construct a context object
	ctx, cancel := newContextTesting(t)
	defer ctx.Done()
	defer cancel()

	// Try the authenticated received
	auth, err := pc.AuthenticatedReceiver(msg, ctx)
	if err != nil {
		t.Errorf("AuthenticatedReceiver() produced an error: %v", err)
	}

	if !auth.IsAuthenticated {
		t.Errorf("Failed: authenticated receiver")
	}

	// Compare the tokens
	if !bytes.Equal(auth.Sender.receptionToken.GetBytes(), tkn[:]) {
		t.Errorf("Tokens do not match! \n\tExptected: "+
			"%+v\n\tReceived: %+v", tkn, auth.Sender.receptionToken)
	}
}

// Error path
func TestProtoComms_AuthenticatedReceiver_BadId(t *testing.T) {
	// Create comm object
	pc := ProtoComms{
		Manager:    newManager(),
		tokens:     token.NewMap(),
		grpcServer: nil,
		privateKey: nil,
	}
	// Create id and token
	testID := id.NewIdFromString("testSender", id.Node, t)
	expectedVal := []byte("testToken")
	tkn := token.Token{}
	copy(tkn[:], expectedVal)

	// Add host
	_, err := pc.AddHost(testID, "", nil, GetDefaultHostParams())
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
		Token:     tkn.Marshal(),
		Message:   nil,
	}

	// Construct a context object
	ctx, cancel := newContextTesting(t)
	defer ctx.Done()
	defer cancel()

	// Try the authenticated received
	a, _ := pc.AuthenticatedReceiver(msg, ctx)

	if a.IsAuthenticated {
		t.Errorf("Expected error path!"+
			"Should not be able to marshal a message with id: %v", badId)
	}

}

// Happy path
func TestProtoComms_GenerateToken(t *testing.T) {
	comm := ProtoComms{
		tokens:     token.NewMap(),
		grpcServer: nil,
		privateKey: nil,
	}
	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	tkn, err := token.Unmarshal(tokenBytes)
	if err != nil {
		t.Errorf("Should be able to unmarshal token: %s", err)
	}

	ok := comm.tokens.Validate(tkn)
	if !ok {
		t.Errorf("Unable to validate token")
	}
}

// Happy path
func TestProtoComms_PackAuthenticatedMessage(t *testing.T) {
	testServerId := id.NewIdFromString("test12345", id.Node, t)
	comm := ProtoComms{
		networkId:  testServerId,
		grpcServer: nil,
		privateKey: nil,
		tokens:     token.NewMap(),
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}

	tkn := token.Token{}
	copy(tkn[:], tokenBytes)

	testId := id.NewIdFromString("test", id.Node, t)

	host, err := NewHost(testId, "test", nil, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tkn)

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
		networkId:  testId,
		grpcServer: nil,
		privateKey: nil,
		tokens:     token.NewMap(),
		Manager:    newManager(),
	}
	err := comm.setPrivateKey(testkeys.LoadFromPath(testkeys.GetNodeKeyPath()))
	if err != nil {
		t.Errorf("Expected to set private key: %+v", err)
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}
	tkn := token.Token{}
	copy(tkn[:], tokenBytes)

	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	host, err := comm.AddHost(testId, "test", pub, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tkn)

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

	if !bytes.Equal(msg.Token, host.transmissionToken.GetBytes()) {
		t.Errorf("Message token doesn't match message's token! "+
			"Expected: %+v"+
			"\n\tReceived: %+v", host.transmissionToken, msg.Token)
	}
}

// Error Path
func TestProtoComms_ValidateToken_BadId(t *testing.T) {
	testId := id.NewIdFromString("test", id.Node, t)
	comm := ProtoComms{
		networkId:  testId,
		grpcServer: nil,
		privateKey: nil,
		tokens:     token.NewMap(),
		Manager:    newManager(),
	}
	err := comm.setPrivateKey(testkeys.LoadFromPath(testkeys.GetNodeKeyPath()))
	if err != nil {
		t.Errorf("Expected to set private key: %+v", err)
	}

	tokenBytes, err := comm.GenerateToken()
	if err != nil || tokenBytes == nil {
		t.Errorf("Unable to generate token: %+v", err)
	}
	tkn := token.Token{}
	copy(tkn[:], tokenBytes)

	pub := testkeys.LoadFromPath(testkeys.GetNodeCertPath())
	host, err := comm.AddHost(testId, "test", pub, GetDefaultHostParams())
	if err != nil {
		t.Errorf("Unable to create host: %+v", err)
	}
	host.transmissionToken.Set(tkn)

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

func TestProtoComms_DisableAuth(t *testing.T) {
	testId := id.NewIdFromString("test", id.Node, t)
	comm := ProtoComms{
		networkId:  testId,
		grpcServer: nil,
		privateKey: nil,
	}

	comm.DisableAuth()

	if !comm.disableAuth {
		t.Error("Auth was not disabled when DisableAuth was called")
	}
}

// newContextTesting constructs a context.Context object on
// the local Unix default domain (UDP) port
func newContextTesting(t *testing.T) (context.Context, context.CancelFunc) {
	protoCtx, cancel := newContext(time.Second)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("udp", "0.0.0.0:53", timeout)
	if err != nil {
		t.Fatalf("Failed to get a conn object in setup: %v", err)
	}
	p := &peer.Peer{
		Addr: conn.RemoteAddr(),
	}

	return peer.NewContext(protoCtx, p), cancel
}
