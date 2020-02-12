////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

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
	"math"
	"net"
	"strings"
	"sync"
	"time"
)

// Proto object containing a gRPC server
type ProtoComms struct {
	// Inherit the Manager object
	Manager

	// The network ID of this comms server
	Id string

	// A map of reverse-authentication tokens
	tokens sync.Map

	// Local network server
	LocalServer *grpc.Server

	// Listening address of the local server
	ListeningAddr string

	// Private key of the local server
	privateKey *rsa.PrivateKey

	//Disables the checking of authentication signatures for testing setups
	disableAuth bool
}

// Starts a ProtoComms object and is used in various initializers
func StartCommServer(id string, localServer string, certPEMblock,
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

	// Ensure the connection is running
	jww.DEBUG.Printf("Attempting to send to host: %s", host)
	if !host.Connected() {
		err = host.connect()
		if err != nil {
			return
		}
	}

	// Number of attempts to negotiate a handshake
	numTries := 1

	// Authentication loop
authenticate:
	numTries--

	// Establish authentication if required
	if host.authenticationRequired() {
		err = host.authenticate(c.clientHandshake)
		if err != nil {
			return
		}
	}
	// Attempt to send to host
	result, err = host.send(f)

	// If failed to authenticate, retry negotiation by jumping to the top of the loop
	if err != nil && strings.Contains(err.Error(), "Failed to authenticate") && numTries > 0 {
		jww.WARN.Printf("Failed to authenticate, %d retries left", numTries)
		goto authenticate
	}

	// Run the send function
	return result, err
}

// Sets up or recovers the Host's connection
// Then runs the given Stream function
func (c *ProtoComms) Stream(host *Host, f func(conn *grpc.ClientConn) (
	interface{}, error)) (client interface{}, err error) {

	// Ensure the connection is running
	jww.DEBUG.Printf("Attempting to send to host: %s", host)
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
		ndfBytes := currentDef.Serialize()
		hash.Write(ndfBytes)
		ndfHash = hash.Sum(nil)
	}
	//Put the hash in a message
	msg := &mixmessages.NDFHash{Hash: ndfHash}

	regHost, ok := c.Manager.GetHost(id.PERMISSIONING)
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
