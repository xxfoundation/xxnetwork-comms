////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for describing and creating connections

package connect

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/crypto/signature/rsa"
	tlsCreds "gitlab.com/elixxir/crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"math"
	"sync"
	"time"
)

// Represents a reverse-authentication token
type Token []byte

// Information used to describe a connection to a host
type Host struct {
	// System-wide ID of the Host
	id string

	// address:Port being connected to
	address string

	// PEM-format TLS Certificate
	certificate []byte

	/* Tokens shared with this Host establishing reverse authentication */

	//  Token used for receiving from this host
	receptionToken Token

	// Token used for sending to this host
	transmissionToken Token

	// Configure the maximum number of connection attempts
	maxRetries int

	// GRPC connection object
	connection *grpc.ClientConn

	// TLS credentials object used to establish the connection
	credentials credentials.TransportCredentials

	// RSA Public Key corresponding to the TLS Certificate
	rsaPublicKey *rsa.PublicKey

	// If set, reverse authentication will be established with this Host
	enableAuth bool

	// Indicates whether dynamic authentication was used for this Host
	// This is useful for determining whether a Host's key was hardcoded
	dynamicHost bool

	// Read/Write Mutex for thread safety
	mux sync.RWMutex
}

// Creates a new Host object
func NewHost(id, address string, cert []byte, disableTimeout,
	enableAuth bool) (host *Host, err error) {

	// Initialize the Host object
	host = &Host{
		id:          id,
		address:     address,
		certificate: cert,
		enableAuth:  enableAuth,
	}

	// Set the max number of retries for establishing a connection
	if disableTimeout {
		host.maxRetries = math.MaxInt32
	} else {
		host.maxRetries = 100
	}

	// Configure the host credentials
	err = host.setCredentials()
	return
}

// Creates a new dynamic-authenticated Host object
func newDynamicHost(id string, publicKey []byte) (host *Host, err error) {

	// Initialize the Host object
	// IMPORTANT: This flag must be set to true for all dynamic Hosts
	//            because the security properties for these Hosts differ
	host = &Host{
		id:          id,
		dynamicHost: true,
	}

	// Create the RSA Public Key object
	host.rsaPublicKey, err = rsa.LoadPublicKeyFromPem(publicKey)
	if err != nil {
		err = errors.Errorf("Error extracting PublicKey: %+v", err)
	}
	return
}

// Simple getter for the dynamicHost value
func (h *Host) IsDynamicHost() bool {
	return h.dynamicHost
}

// Simple getter for the public key
func (h *Host) GetPubKey() *rsa.PublicKey {
	return h.rsaPublicKey
}

// Connected checks if the given Host's connection is alive
func (h *Host) Connected() bool {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.isAlive()
}

// GetId  returns the id of the host
func (h *Host) GetId() string {
	return h.id
}

// Disconnect closes a the Host connection under the write lock
func (h *Host) Disconnect() {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.disconnect()
	h.receptionToken = nil
	h.transmissionToken = nil
}

// send checks that the host has a connection and sends if it does.
// Operates under the host's read lock.
func (h *Host) send(f func(conn *grpc.ClientConn) (*any.Any,
	error)) (*any.Any, error) {

	h.mux.RLock()
	defer h.mux.RUnlock()

	if !h.isAlive() {
		return nil, errors.New("Could not send, connection is not alive")
	}

	a, err := f(h.connection)
	return a, err
}

// stream checks that the host has a connection and streams if it does.
// Operates under the host's read lock.
func (h *Host) stream(f func(conn *grpc.ClientConn) (
	interface{}, error)) (interface{}, error) {

	h.mux.RLock()
	defer h.mux.RUnlock()

	if !h.isAlive() {
		return nil, errors.New("Could not stream, connection is not alive")
	}

	a, err := f(h.connection)
	return a, err
}

// connect attempts to connect to the host if it does not have a valid connection
func (h *Host) connect() error {
	h.mux.Lock()
	defer h.mux.Unlock()

	//checks if the connection is active and skips reconnecting if it is
	if h.isAlive() {
		return nil
	}

	//connect to remote
	if err := h.connectHelper(); err != nil {
		return err
	}

	return nil
}

// authenticationRequired Checks if new authentication is required with
// the remote
func (h *Host) authenticationRequired() bool {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.enableAuth && h.transmissionToken == nil
}

// Checks if the given Host's connection is alive
func (h *Host) authenticate(handshake func(host *Host) error) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	return handshake(h)
}

// isAlive returns true if the connection is non-nil and alive
func (h *Host) isAlive() bool {
	if h.connection == nil {
		return false
	}
	state := h.connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

// disconnect closes a the Host connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (h *Host) disconnect() {
	// its possible to close a host which never sent so it never made a
	// connection. In that case, we should not close a connection which does not
	// exist
	if h.connection != nil {
		err := h.connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v",
				h.address, errors.New(err.Error()))
		}
	}
}

// connectHelper creates a connection while not under a write lock.
// undefined behavior if the caller has not taken the write lock
func (h *Host) connectHelper() (err error) {

	// Configure TLS options
	var securityDial grpc.DialOption
	if h.credentials != nil {
		// Create the gRPC client with TLS
		securityDial = grpc.WithTransportCredentials(h.credentials)
	} else {
		// Create the gRPC client without TLS
		jww.WARN.Printf("Connecting to %v without TLS!", h.address)
		securityDial = grpc.WithInsecure()
	}

	jww.DEBUG.Printf("Attempting to establish connection to %s using"+
		" credentials: %+v", h.address, securityDial)

	// Attempt to establish a new connection
	for numRetries := 0; numRetries < h.maxRetries && !h.isAlive(); numRetries++ {

		jww.INFO.Printf("Connecting to %+v. Attempt number %+v of %+v",
			h.address, numRetries, h.maxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2 * (numRetries/16 + 1)
		if backoffTime > 15 {
			backoffTime = 15
		}
		ctx, cancel := ConnectionContext(time.Duration(backoffTime))

		// Create the connection
		h.connection, err = grpc.DialContext(ctx, h.address, securityDial,
			grpc.WithBlock(), grpc.WithBackoffMaxDelay(time.Minute*5))
		if err != nil {
			jww.ERROR.Printf("Attempt number %+v to connect to %s failed: %+v\n",
				numRetries, h.address, errors.New(err.Error()))
		}
		cancel()
	}

	// Verify that the connection was established successfully
	if !h.isAlive() {
		return errors.New(fmt.Sprintf(
			"Last try to connect to %s failed. Giving up", h.address))
	}

	// Add the successful connection to the Manager
	jww.INFO.Printf("Successfully connected to %v", h.address)
	return
}

// setCredentials sets TransportCredentials and RSA PublicKey objects
// using a PEM-encoded TLS Certificate
func (h *Host) setCredentials() error {

	// If no TLS Certificate specified, print a warning and do nothing
	if h.certificate == nil || len(h.certificate) == 0 {
		jww.WARN.Printf("No TLS Certificate specified!")
		return nil
	}

	// Obtain the DNS name included with the certificate
	dnsName := ""
	cert, err := tlsCreds.LoadCertificate(string(h.certificate))
	if err != nil {
		return errors.Errorf("Error forming transportCredentials: %+v", err)
	}
	if len(cert.DNSNames) > 0 {
		dnsName = cert.DNSNames[0]
	}

	// Create the TLS Credentials object
	h.credentials, err = tlsCreds.NewCredentialsFromPEM(string(h.certificate),
		dnsName)
	if err != nil {
		return errors.Errorf("Error forming transportCredentials: %+v", err)
	}

	// Create the RSA Public Key object
	h.rsaPublicKey, err = tlsCreds.NewPublicKeyFromPEM(h.certificate)
	if err != nil {
		err = errors.Errorf("Error extracting PublicKey: %+v", err)
	}

	return err
}

// Stringer interface for connection
func (h *Host) String() string {
	addr := h.address
	actualConnection := h.connection
	creds := h.credentials

	var state connectivity.State
	if actualConnection != nil {
		state = actualConnection.GetState()
	}

	serverName := "<nil>"
	protocolVersion := "<nil>"
	securityVersion := "<nil>"
	securityProtocol := "<nil>"
	if creds != nil {
		serverName = creds.Info().ServerName
		securityVersion = creds.Info().SecurityVersion
		protocolVersion = creds.Info().ProtocolVersion
		securityProtocol = creds.Info().SecurityProtocol
	}
	cStrt := len("-----BEGIN CERTIFICATE----- ") // Skip this part
	return fmt.Sprintf(
		"ID: %v\tAddr: %v\tCertificate: %s...\tTransmission Token: %v"+
			"\tReception Token: %+v \tEnableAuth: %v"+
			"\tMaxRetries: %v\tConnState: %v"+
			"\tTLS ServerName: %v\tTLS ProtocolVersion: %v\t"+
			"TLS SecurityVersion: %v\tTLS SecurityProtocol: %v\n",
		h.id, addr, h.certificate[cStrt:cStrt+20], h.transmissionToken,
		h.receptionToken, h.enableAuth, h.maxRetries, state,
		serverName, protocolVersion, securityVersion, securityProtocol)
}
