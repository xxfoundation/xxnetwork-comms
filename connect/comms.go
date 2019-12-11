////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

//

package connect

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"google.golang.org/grpc"
	"time"
)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	Manager

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
