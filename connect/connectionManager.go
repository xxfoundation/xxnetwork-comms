////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for connecting to gateways and servers

package connect

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/signature/rsa"
	tlsCreds "gitlab.com/elixxir/crypto/tls"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"math"
	"sort"
	"sync"
	"time"
)

// Stores a connection and its associated information
type connection struct {
	// Address:Port being connected to
	Address string

	// GRPC connection object
	Connection *grpc.ClientConn

	// Credentials object used to establish the connection
	Creds credentials.TransportCredentials

	// RSA Public Key corresponding to the TLS Certificate
	RsaPublicKey *rsa.PublicKey
}

// Information used to describe a connection
type ConnectionInfo struct {
	// ID used to identify the connection
	Id fmt.Stringer

	// Address:Port being connected to
	Address string

	// PEM-format TLS Certificate
	Cert []byte

	// Indicates whether connection timeout should be disabled
	DisableTimeout bool
}

type ConnectionManager struct {
	// A map of string IDs to open connections
	connections map[string]*connection

	// Local server RSA Private Key
	privateKey *rsa.PrivateKey

	connectionsLock sync.Mutex
	maxRetries      int64
}

// Default maximum number of retries
const DefaultMaxRetries = 100

// Set private key to data to a PEM block
func (m *ConnectionManager) SetPrivateKey(data []byte) error {
	key, err := rsa.LoadPrivateKeyFromPem(data)
	if err != nil {
		s := fmt.Sprintf("Failed to form private key file from data at %s: %+v", data, err)
		return errors.New(s)
	}

	m.privateKey = key
	return nil
}

// Get connection manager's private key
func (m *ConnectionManager) GetPrivateKey() *rsa.PrivateKey {
	return m.privateKey
}

func (m *ConnectionManager) GetConnection(id string) *connection {
	return m.connections[id]
}

func (m *ConnectionManager) SetMaxRetries(mr int64) {
	m.maxRetries = mr
}

// Creates TransportCredentials and RSA PublicKey objects
// using a PEM-encoded TLS Certificate
func createCredentials(connInfo *ConnectionInfo) (credentials.
	TransportCredentials, *rsa.PublicKey, error) {

	// If no TLS Certificate specified, print a warning and do nothing
	if connInfo.Cert == nil || len(connInfo.Cert) == 0 {
		jww.WARN.Printf("No TLS Certificate specified!")
		return nil, nil, nil
	}

	// Obtain the DNS name included with the certificate
	dnsName := ""
	cert, err := tlsCreds.LoadCertificate(string(connInfo.Cert))
	if err != nil {
		s := fmt.Sprintf("Error forming transportCredentials: %+v", err)
		return nil, nil, errors.New(s)
	}
	if len(cert.DNSNames) > 0 {
		dnsName = cert.DNSNames[0]
	}

	// Create the TLS Credentials object
	tlsCredentials, err := tlsCreds.NewCredentialsFromPEM(string(connInfo.Cert),
		dnsName)
	if err != nil {
		s := fmt.Sprintf("Error forming transportCredentials: %+v", err)
		return nil, nil, errors.New(s)
	}

	// Create the RSA Public Key object
	publicKey, err := tlsCreds.NewPublicKeyFromPEM(connInfo.Cert)
	if err != nil {
		s := fmt.Sprintf("Error extracting PublicKey: %+v", err)
		return nil, nil, errors.New(s)
	}

	return tlsCredentials, publicKey, nil
}

// ConnectToRemote connects to a remote server at address addr with the passed
// cert. The connection is stored locally at the passed id.  that ID can be
// used to identify the private keys of the sender of incoming messages so
// it must be the same as used across the network.
// TODO: Deprecated. Create connections automatically if they do not exist
func (m *ConnectionManager) ConnectToRemote(connInfo *ConnectionInfo) error {
	return m.connect(connInfo)
}

func (m *ConnectionManager) GetRegistrationConnection(connInfo *ConnectionInfo) pb.
	RegistrationClient {
	conn := m.get(connInfo)
	if !isConnectionGood(conn) {
		jww.WARN.Printf("Bad Registration connection state, "+
			"reconnecting: %v",
			m.connections[connInfo.Id.String()])
		resetConnection(conn)
	}
	return pb.NewRegistrationClient(conn)
}

func (m *ConnectionManager) GetGatewayConnection(connInfo *ConnectionInfo) pb.
	GatewayClient {
	conn := m.get(connInfo)
	if !isConnectionGood(conn) {
		jww.WARN.Printf("Bad Gateway connection state, "+
			"reconnecting: %v",
			m.connections[connInfo.Id.String()])
		resetConnection(conn)
	}
	return pb.NewGatewayClient(conn)
}

func (m *ConnectionManager) GetNodeConnection(connInfo *ConnectionInfo) pb.
	NodeClient {
	conn := m.get(connInfo)
	if !isConnectionGood(conn) {
		jww.WARN.Printf("Bad Node connection state, reconnecting: %v",
			m.connections[connInfo.Id.String()])
		resetConnection(conn)
	}
	return pb.NewNodeClient(conn)
}

// Attempts to reconnect to the remote host
func resetConnection(connection *grpc.ClientConn) {
	// NOTE: This is currently experimental, but claims to immediately
	//       reconnect. We wrap this so we can fix/change later...
	// https://godoc.org/google.golang.org/grpc#ClientConn.ResetConnectBackoff
	connection.ResetConnectBackoff()
}

// Returns true if the connection is non-nil and alive
func isConnectionGood(connection *grpc.ClientConn) bool {
	if connection == nil {
		return false
	}
	state := connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready
}

// Get creates an existing connection
func (m *ConnectionManager) get(connInfo *ConnectionInfo) *grpc.ClientConn {
	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()

	conn, ok := m.connections[connInfo.Id.String()]
	if !ok {
		// TODO: Create the connection if it does not exist
		jww.FATAL.Panicf("No connection exists for ID %s",
			connInfo.Id.String())
	}
	return conn.Connection
}

// Connect creates a connection
func (m *ConnectionManager) connect(connInfo *ConnectionInfo) error {
	var clientConnection *grpc.ClientConn
	var securityDial grpc.DialOption

	// Obtain credentials
	tlsCredentials, publicKey, err := createCredentials(connInfo)
	if err != nil {
		return err
	}

	if tlsCredentials != nil {
		// Create the gRPC client with TLS
		securityDial = grpc.WithTransportCredentials(tlsCredentials)
	} else {
		// Create the gRPC client without TLS
		jww.WARN.Printf("Connecting to %v without TLS!", connInfo.Address)
		securityDial = grpc.WithInsecure()
	}

	// Initialize the connections map if it hasn't already
	if m.connections == nil {
		m.connections = make(map[string]*connection)
	}

	// Set the max number of retries for establishing a connection
	var maxRetries int64
	if connInfo.DisableTimeout {
		maxRetries = math.MaxInt64
	} else {
		if m.maxRetries == 0 {
			maxRetries = DefaultMaxRetries
		} else {
			maxRetries = m.maxRetries
		}
	}

	// Attempt to create a new connection
	for numRetries := int64(0); numRetries < maxRetries && !isConnectionGood(clientConnection); numRetries++ {

		jww.INFO.Printf("Connecting to address %+v. Attempt number %+v of %+v",
			connInfo.Address, numRetries, maxRetries)

		// If timeout is enabled, the max wait time becomes
		// ~14 seconds (with maxRetries=100)
		backoffTime := 2 * (numRetries/16 + 1)
		if backoffTime > 15 {
			backoffTime = 15
		}
		ctx, cancel := ConnectionContext(time.Duration(backoffTime))

		// Create the connection
		clientConnection, err = grpc.DialContext(ctx, connInfo.Address,
			securityDial,
			grpc.WithBlock(),
			grpc.WithBackoffMaxDelay(time.Minute*5))
		if err != nil {
			jww.ERROR.Printf("Attempt number %+v to connect to %s failed: %+v\n",
				numRetries, connInfo.Address, errors.New(err.Error()))
		}
		cancel()
	}

	// Verify that the connection was established successfully
	if !isConnectionGood(clientConnection) {
		return errors.New(fmt.Sprintf("Last try to connect to %s failed. Giving up",
			connInfo.Address))
	}

	// Add the successful connection to the ConnectionManager
	jww.INFO.Printf("Successfully connected to %s at %v",
		connInfo.Id, connInfo.Address)
	m.connectionsLock.Lock()
	m.connections[connInfo.Id.String()] = &connection{
		Address:      connInfo.Address,
		Connection:   clientConnection,
		Creds:        tlsCredentials,
		RsaPublicKey: publicKey,
	}
	m.connectionsLock.Unlock()
	return nil
}

// Disconnect closes client connections and removes it from ConnectionManager
func (m *ConnectionManager) Disconnect(id string) {
	m.connectionsLock.Lock()
	connection, present := m.connections[id]
	if present {
		err := connection.Connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %+v", id,
				errors.New(err.Error()))
		}
		delete(m.connections, id)
	}
	m.connectionsLock.Unlock()
}

// DisconnectAll closes all client connections
// and removes them from ConnectionManager
func (m *ConnectionManager) DisconnectAll() {
	for connId := range m.connections {
		m.Disconnect(connId)
	}
}

// Implements Stringer for debug printing
func (m *ConnectionManager) String() string {
	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()

	// Sort connection IDs to print in a consistent order
	keys := make([]string, len(m.connections))
	i := 0
	for key := range m.connections {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Print each connection's information
	var result bytes.Buffer
	for _, key := range keys {
		// Populate fields without ever dereferencing nil
		connection := m.connections[key]
		if connection != nil {
			result.WriteString(fmt.Sprintf(
				"[%v] %s",
				key, connection))
		}
	}

	return result.String()
}

// Stringer interface for connection
func (ci *connection) String() string {
	addr := ci.Address
	actualConnection := ci.Connection
	creds := ci.Creds

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
	return fmt.Sprintf(
		"Addr: %v\tState: %v\tTLS ServerName: %v\t"+
			"TLS ProtocolVersion: %v\tTLS SecurityVersion: %v\t"+
			"TLS SecurityProtocol: %v\n",
		addr, state, serverName, protocolVersion,
		securityVersion, securityProtocol)
}

// TimeoutContext is a context with the default timeout
func ConnectionContext(seconds time.Duration) (context.Context, context.CancelFunc) {
	waitingPeriod := seconds * time.Second
	jww.DEBUG.Printf("Timing out in: %s", waitingPeriod)
	ctx, cancel := context.WithTimeout(context.Background(),
		waitingPeriod)
	return ctx, cancel

}

// DefaultContexts creates a context object with the default context
// for all client messages. This is primarily used to set the default
// timeout for all clients
func MessagingContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(),
		2*time.Minute)
	return ctx, cancel
}

// StreamingContext creates a context object with the default context
// for all client streaming messages. This is primarily used to
// allow a cancel option for clients and is suitable for unary streaming.
func StreamingContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel
}
