////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/rand"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/nonce"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	Manager

	// A map of reverse-authentication tokens
	tokens sync.Map

	// Local network server
	LocalServer *grpc.Server

	// Listening address of the local server
	ListeningAddr string

	// Private key of the local server
	privateKey *rsa.PrivateKey
}

// Performs a graceful shutdown of the local server
func (c *ProtoComms) Shutdown() {
	c.DisconnectAll()
	c.LocalServer.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// Stringer method
func (c *ProtoComms) String() string {
	return c.ListeningAddr
}

// Setter for local server's private key
func (c *ProtoComms) SetPrivateKey(data []byte) error {
	key, err := rsa.LoadPrivateKeyFromPem(data)
	if err != nil {
		return errors.Errorf("Failed to form private key file from data at %s: %+v", data, err)
	}

	c.privateKey = key
	return nil
}

// Getter for local server's private key
func (c *ProtoComms) GetPrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

// Convert any message type into a authenticated message
func (c *ProtoComms) Authenticate(msg proto.Message, host *Host,
	enableSignature bool) (*pb.AuthenticatedMessage, error) {

	// Marshall the provided message into an Any type
	anyMsg, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Build the authenticated message
	authMsg := &pb.AuthenticatedMessage{
		ID:        host.id,
		Signature: nil,
		Token:     host.token,
		Message:   anyMsg,
	}

	// If signature is enabled, sign the message and add to payload
	if enableSignature {
		authMsg.Signature, err = c.signMessage(anyMsg)
		if err != nil {
			return nil, err
		}
	}

	return authMsg, nil
}

// Generates a new token and adds it to internal state
func (c *ProtoComms) GenerateToken() ([]byte, error) {
	token, err := nonce.NewNonce(nonce.RegistrationTTL)
	if err != nil {
		return nil, err
	}

	c.tokens.Store(token.Bytes(), token)
	return token.Bytes(), nil
}

// Validates an authenticated message using internal state
func (c *ProtoComms) ValidateToken(msg *pb.AuthenticatedMessage) error {

	// Verify the token was assigned
	token, ok := c.tokens.Load(msg.Token)
	if !ok {
		return errors.Errorf("Unable to locate token: %+v", msg.Token)
	}

	// Verify the token is not expired
	if !token.(*nonce.Nonce).IsValid() {
		return errors.Errorf("Invalid or expired token: %+v", msg.Token)
	}

	// Verify the Host exists for the provided ID
	host, ok := c.GetHost(msg.ID)
	if !ok {
		return errors.Errorf("Invalid token for host ID: %+v", msg.ID)
	}

	// Verify the token signature
	if err := c.verifyMessage(msg, host); err != nil {
		return errors.Errorf("Invalid token signature: %+v", err)
	}

	// Token has been validated and can be safely stored
	host.SetToken(msg.Token)
	return nil
}

// Takes a generic-type message, returns the signature
// The message is signed with the ProtoComms RSA PrivateKey
func (c *ProtoComms) signMessage(anyMessage *any.Any) ([]byte, error) {
	// Hash the message data
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	data := []byte(anyMessage.String())
	hashed := hash.Sum(data)[len(data):]

	// Obtain the private key
	key := c.GetPrivateKey()
	if key == nil {
		return nil, errors.Errorf("Cannot sign message: No private key")
	}

	// Sign the message and return the signature
	signature, err := rsa.Sign(rand.Reader, key, options.Hash, hashed, nil)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return signature, nil
}

// Takes an AuthenticatedMessage and a Host, verifies the signature
// using Host public key, returning an error if invalid
func (c *ProtoComms) verifyMessage(msg *pb.AuthenticatedMessage, host *Host) error {

	// Get hashed data of the message
	options := rsa.NewDefaultOptions()
	hash := options.Hash.New()
	s := msg.Message.String()
	data := []byte(s)
	hashed := hash.Sum(data)[len(data):]

	// Verify signature of message using host public key
	err := rsa.Verify(host.rsaPublicKey, options.Hash, hashed, msg.Signature, nil)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
