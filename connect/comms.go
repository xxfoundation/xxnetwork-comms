////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
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
func (c *ProtoComms) Authenticate(msg proto.Message, host *Host) (
	*pb.AuthenticatedMessage, error) {

	// Marshall the provided message into an Any type
	anyMsg, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Build the authenticated message
	// NOTE: Signature remains unused until needed
	return &pb.AuthenticatedMessage{
		ID:        host.id,
		Signature: nil,
		Token:     host.token,
		Message:   anyMsg,
	}, nil
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
	token, ok := c.tokens.Load(msg.Token)
	if !ok {
	}

	if !token.(*nonce.Nonce).IsValid() {
	}

	host, ok := c.GetHost(msg.ID)
	if !ok {

	}

	host.SetToken(msg.Token)
	return nil
}
