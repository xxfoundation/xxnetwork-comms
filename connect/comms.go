///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Handles the basic top-level comms object used across all packages

package connect

import (
	"crypto/sha256"
	"crypto/tls"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
	"gitlab.com/elixxir/primitives/id"
	"gitlab.com/elixxir/primitives/ndf"
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
	// Close conn 2 seconds afer ping
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

	numConnects, numAuths, lastEvent := 0, 0, 0

connect:
	// Ensure the connection is running
	if !host.Connected() {
		host.transmissionToken = nil
		//do not attempt to connect again if multiple attempts have been made
		if numConnects == maxConnects {
			return nil, errors.WithMessage(err, "Maximum number of connects attempted")
		}

		//denote that a connection is being tried
		lastEvent = con

		//attempt to make the connection
		jww.INFO.Printf("Host %s disconnected, attempting to reconnect...", host.id.String())
		err = host.connect()
		//if connection cannot be made, do not retry
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to connect")
		}

		//denote the connection attempt
		numConnects++
	}

authorize:
	// Establish authentication if required
	if host.authenticationRequired() && host.transmissionToken == nil {
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
			//otherwise, return the error
			return nil, errors.New("Failed to authenticate")
		}

		//denote the authorization attempt
		numAuths++
	}

	//denote that a send is being tried
	lastEvent = send
	// Attempt to send to host
	result, err = host.send(f)
	// If failed to authenticate, retry negotiation by jumping to the top of the loop
	if err != nil {
		//if failure of connection, retry connection
		if isConnError(err) {
			jww.INFO.Printf("Failed send due to connection issue: %s", err)
			goto connect
		}

		// Handle resetting authentication
		if strings.Contains(err.Error(), AuthError(host.id).Error()) {
			jww.INFO.Printf("Failed send due to auth error, retrying authentication: %s", err.Error())
			host.transmissionToken = nil
			goto authorize
		}

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

// RequestNdf is used to Request an ndf from permissioning
// Used by gateway, client, nodes and gateways
func (c *ProtoComms) RequestNdf(host *Host,
	message *mixmessages.NDFHash) (*mixmessages.NDF, error) {

	// Create the Send Function
	f := func(conn *grpc.ClientConn) (*any.Any, error) {
		// Set up the context
		ctx, cancel := MessagingContext()
		defer cancel()

		authMsg, err := c.PackAuthenticatedMessage(message, host, false)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		// Send the message
		resultMsg, err := mixmessages.NewRegistrationClient(
			conn).PollNdf(ctx, authMsg)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return ptypes.MarshalAny(resultMsg)
	}

	// Execute the Send function
	jww.DEBUG.Printf("Sending Request Ndf message: %+v", message)
	resultMsg, err := c.Send(host, f)
	if err != nil {
		return nil, err
	}

	result := &mixmessages.NDF{}
	return result, ptypes.UnmarshalAny(resultMsg, result)

}

// RetrieveNdf, attempts to connect to the permissioning server to retrieve the latest ndf for the notifications bot
func (c *ProtoComms) RetrieveNdf(currentDef *ndf.NetworkDefinition) (*ndf.NetworkDefinition, error) {
	//Hash the notifications bot ndf for comparison with registration's ndf
	var ndfHash []byte
	// If the ndf passed not nil, serialize and hash it
	if currentDef != nil {
		//Hash the notifications bot ndf for comparison with registration's ndf
		hash := sha256.New()
		ndfBytes, err := currentDef.Marshal()
		if err != nil {
			return nil, err
		}
		hash.Write(ndfBytes)
		ndfHash = hash.Sum(nil)
	}
	//Put the hash in a message
	msg := &mixmessages.NDFHash{Hash: ndfHash}

	regHost, ok := c.Manager.GetHost(&id.Permissioning)
	if !ok {
		return nil, errors.New("Failed to find permissioning host")
	}

	//Send the hash to registration
	response, err := c.RequestNdf(regHost, msg)

	// Keep going until we get a grpc error or we get an ndf
	for err != nil {
		// If there is an unexpected error
		if !strings.Contains(err.Error(), ndf.NO_NDF) {
			// If it is not an issue with no ndf, return the error up the stack
			errMsg := errors.Errorf("Failed to get ndf from permissioning: %v", err)
			return nil, errMsg
		}

		// If the error is that the permissioning server is not ready, ask again
		jww.WARN.Println("Failed to get an ndf, possibly not ready yet. Retying now...")
		time.Sleep(250 * time.Millisecond)
		response, err = c.RequestNdf(regHost, msg)

	}

	//If there was no error and the response is nil, client's ndf is up-to-date
	if response == nil || response.Ndf == nil {
		jww.DEBUG.Printf("Our NDF is up-to-date")
		return nil, nil
	}

	jww.INFO.Printf("Remote NDF: %s", string(response.Ndf))

	//Otherwise pull the ndf out of the response
	updatedNdf, _, err := ndf.DecodeNDF(string(response.Ndf))
	if err != nil {
		//If there was an error decoding ndf
		errMsg := errors.Errorf("Failed to decode response to ndf: %v", err)
		return nil, errMsg
	}
	return updatedNdf, nil
}
