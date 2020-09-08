///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/tls"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"math"
	"net"
	"strings"
	"sync"
	"time"
)

// TODO: Set these via config

// KaOpts are Keepalive options for servers
var KaOpts = keepalive.ServerParameters{
	// Idle for at most 5s
	MaxConnectionIdle: 5 * time.Second,
	// Reset after an hour
	MaxConnectionAge: 1 * time.Hour,
	// w/ 1m grace shutdown
	MaxConnectionAgeGrace: 1 * time.Minute,
	// ping if no activity after 1s
	Time: 1 * time.Second,
	// Close conn 2 seconds after ping
	Timeout: 2 * time.Second,
}

// KaEnforcement are keepalive enforcement options for servers
var KaEnforcement = keepalive.EnforcementPolicy{
	// Client should wait at least 250ms
	MinTime: 250 * time.Millisecond,
	// Doing KA on non-streams is OK
	PermitWithoutStream: true,
}

// MaxConcurrentStreams is the number of server-side streams to allow open
var MaxConcurrentStreams = uint32(250000)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	Manager

	// The network ID of this comms server
	Id *id.ID

	// Private key of the local comms instance
	privateKey *rsa.PrivateKey

	// Disables the checking of authentication signatures for testing setups
	disableAuth bool

	// SERVER-ONLY FIELDS ------------------------------------------------------

	// A map of reverse-authentication tokens
	tokens sync.Map

	// Local network server
	LocalServer *grpc.Server

	// Listening address of the local server
	ListeningAddr string

	// CLIENT-ONLY FIELDS ------------------------------------------------------

	// Used to store the public key used for generating Client Id
	pubKeyPem []byte

	// Used to store the salt used for generating Client Id
	salt []byte

	// -------------------------------------------------------------------------
}

// Creates a ProtoComms client-type object to be used in various initializers
func CreateCommClient(id *id.ID, pubKeyPem, privKeyPem,
	salt []byte) (*ProtoComms, error) {
	// Build the ProtoComms object
	pc := &ProtoComms{
		Id:        id,
		pubKeyPem: pubKeyPem,
		salt:      salt,
	}

	// Set the private key if specified
	if privKeyPem != nil {
		err := pc.setPrivateKey(privKeyPem)
		if err != nil {
			return nil, errors.Errorf("Could not set private key: %+v", err)
		}
	}
	return pc, nil
}

// Creates a ProtoComms server-type object to be used in various initializers
func StartCommServer(id *id.ID, localServer string, certPEMblock,
	keyPEMblock []byte) (*ProtoComms, net.Listener, error) {

	// Build the ProtoComms object
	pc := &ProtoComms{
		Id:            id,
		ListeningAddr: localServer,
	}

listen:
	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			jww.WARN.Printf("Could not listen on %s, is port in use? waiting 30s: %s", localServer, err.Error())
			time.Sleep(30 * time.Second)
			goto listen
		}
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
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
	} else {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		pc.LocalServer = grpc.NewServer(
			grpc.MaxConcurrentStreams(MaxConcurrentStreams),
			grpc.MaxRecvMsgSize(math.MaxInt32),
			grpc.KeepaliveParams(KaOpts),
			grpc.KeepaliveEnforcementPolicy(KaEnforcement))
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

const (
	//numerical designators for 3 operations of send, used to ensure the same
	//operation isn't repeated
	con  = 1
	auth = 2
	send = 3
	//maximum number of connection attempts. should be 1 because a connection
	//is unlikely to be successful after a failure
	maxConnects = 1
	//maximum number of auth attempts. should be 3 because while a connection
	//is unlikely to be successful after a failure, it is possible to attempt
	//one on a dead connection and then need to do it again after the connection
	//is resuscitated
	maxAuths = 2
)

// Sets up or recovers the Host's connection
// Then runs the given Send function
func (c *ProtoComms) Send(host *Host, f func(conn *grpc.ClientConn) (*any.Any,
	error)) (result *any.Any, err error) {

	if host.GetAddress() == "" {
		return nil, errors.New("Host address is blank, host might be receive only.")
	}

	numConnects, numAuths, lastEvent := 0, 0, 0
	host.sendLock.Lock()
connect:
	// Ensure the connection is running
	if !host.Connected() {
		host.transmissionToken.SetToken(nil)
		//do not attempt to connect again if multiple attempts have been made
		if numConnects == maxConnects {
			host.sendLock.Unlock()
			return nil, errors.WithMessage(err, "Maximum number of connects attempted")
		}

		//denote that a connection is being tried
		lastEvent = con

		//attempt to make the connection
		jww.INFO.Printf("Host %s not connected, attempting to connect...", host.id.String())
		err = host.connect()
		//if connection cannot be made, do not retry
		if err != nil {
			host.sendLock.Unlock()
			return nil, errors.WithMessage(err, "Failed to connect")
		}

		//denote the connection attempt
		numConnects++
	}

authorize:
	// Establish authentication if required
	if host.authenticationRequired() && host.transmissionToken.GetToken() == nil {
		//do not attempt to connect again if multiple attempts have been made
		if numAuths == maxAuths {
			return nil, errors.New("Maximum number of authorizations attempted")
		}

		//do not try multiple auths in a row
		if lastEvent == auth {
			return nil, errors.New("Cannot attempt to authorize with host multiple times in a row")
		}

		//denote that an auth is being tried
		lastEvent = auth

		jww.INFO.Printf("Attempting to establish authentication with host %s", host.id.String())
		err = host.authenticate(c.clientHandshake)
		if err != nil {
			//if failure of connection, retry connection
			if isConnError(err) {
				jww.INFO.Printf("Failed to auth due to connection issue: %s", err)
				goto connect
			}
			host.sendLock.Unlock()
			//otherwise, return the error
			return nil, errors.New("Failed to authenticate")
		}

		//denote the authorization attempt
		numAuths++
	}

	//denote that a send is being tried
	lastEvent = send
	// Attempt to send to host
	host.sendLock.Unlock()
	result, err = host.send(f)
	// If failed to authenticate, retry negotiation by jumping to the top of the loop
	if err != nil {
		//if failure of connection, retry connection
		if isConnError(err) {
			host.sendLock.Lock()
			jww.INFO.Printf("Failed send due to connection issue: %s", err)
			goto connect
		}

		// Handle resetting authentication
		if strings.Contains(err.Error(), AuthError(host.id).Error()) {
			jww.INFO.Printf("Failed send due to auth error, retrying authentication: %s", err.Error())
			host.transmissionToken.SetToken(nil)
			host.sendLock.Lock()
			goto authorize
		}
		host.sendLock.Unlock()
		// otherwise, return the error
		return nil, errors.WithMessage(err, "Failed to send")
	}

	return result, err
}

// returns true if the connection error is one of the connection errors which
// should be retried
func isConnError(err error) bool {
	return strings.Contains(err.Error(), "context deadline exceeded") ||
		strings.Contains(err.Error(), "connection refused")
}

// Sets up or recovers the Host's connection
// Then runs the given Stream function
func (c *ProtoComms) Stream(host *Host, f func(conn *grpc.ClientConn) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.TRACE.Printf("Attempting to send to host: %s", host)
	if !host.Connected() {
		err = host.connect()
		if err != nil {
			return
		}
	}

	//establish authentication if required
	if host.authenticationRequired() {
		err = host.authenticate(c.clientHandshake)
		if err != nil {
			return
		}
	}

	// Run the send function
	return host.stream(f)
}
