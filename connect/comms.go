////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/tls"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math"
	"net"
	"sync"
	"testing"
	"time"
)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	Manager

	// The network ID of this comms server
	id string

	// A map of reverse-authentication tokens
	tokens sync.Map

	// Local network server
	LocalServer *grpc.Server

	// Listening address of the local server
	ListeningAddr string

	// Private key of the local server
	privateKey *rsa.PrivateKey

	//Disables the checking of authentication signatures for testing setups
	fakeAuth bool
}

// Starts a ProtoComms object and is used in various initializers
func StartCommServer(id string, localServer string, certPEMblock,
	keyPEMblock []byte) (*ProtoComms, net.Listener, error) {

	// Build the ProtoComms object
	pc := &ProtoComms{
		id:            id,
		ListeningAddr: localServer,
	}

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		return nil, nil, errors.New(err.Error())
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {

		// Create the TLS certificate
		x509cert, err := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err != nil {
			return nil, nil, errors.Errorf("Could not load TLS keys: %+v", err)
		}

		// Set the private key
		err = pc.setPrivateKey(keyPEMblock)
		if err != nil {
			return nil, nil, errors.Errorf("Could not set private key: %+v", err)
		}

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		creds := credentials.NewServerTLSFromCert(&x509cert)
		pc.LocalServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))

	} else {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		pc.LocalServer = grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))
	}

	return pc, lis, nil
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
func (c *ProtoComms) setPrivateKey(data []byte) error {
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

// Sets up or recovers the Host's connection
// Then runs the given Send function
func (c *ProtoComms) Send(host *Host, f func(conn *grpc.ClientConn) (*any.Any,
	error)) (result *any.Any, err error) {

	// Handle thread safety
	host.mux.RLock()

	// Ensure the connection is running
	jww.DEBUG.Printf("Attempting to send to host: %s", host)
	if err = host.validateConnection(); err != nil {
		host.mux.RUnlock()
		return
	}

	// If authentication is enabled and not yet configured, perform handshake
	if host.enableAuth && host.token == nil {
		host.mux.RUnlock()
		if err = c.clientHandshake(host); err != nil {
			return
		}
		host.mux.RLock()
	}

	// Run the send function
	result, err = f(host.connection)
	host.mux.RUnlock()
	return result, err
}

// Sets up or recovers the Host's connection
// Then runs the given Stream function
func (c *ProtoComms) Stream(host *Host, f func(conn *grpc.ClientConn) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.DEBUG.Printf("Attempting to stream to host: %s", host)
	if err = host.validateConnection(); err != nil {
		return
	}

	// Run the stream function
	return f(host.connection)
}

// Makes the authentication code skip signing and signature verification if the
// set.  Can only be set while in a testing structure.  Is not thread safe.
func (c *ProtoComms) FakeAuth(t interface{}) {
	switch t.(type) {
	case testing.B:
	case testing.T:
	case testing.M:
	default:
		jww.FATAL.Panicf("Cannot ")
	}
	c.fakeAuth = true
}
