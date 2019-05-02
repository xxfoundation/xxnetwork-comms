////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for connecting to gateways and servers

package connect

import (
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/utils"
)

// Stores information used to connect to a server
type ConnectionInfo struct {
	Address string
	// You can also get the server name from Creds if you need it
	Creds      credentials.TransportCredentials
	Connection *grpc.ClientConn
}

// Create credentials from a PEM string
// Intended for mobile clients that can't reasonably use a file
func NewCredentialsFromPEM(serverName string, certificate string) credentials.TransportCredentials {
	// Create cert pool
	pool := x509.NewCertPool()
	// Append the cert string
	if !pool.AppendCertsFromPEM([]byte(certificate)) {
		jww.FATAL.Panicf("Failed to parse certificate string!")
	}
	// Generate credentials from pool
	return credentials.NewClientTLSFromCert(pool, serverName)
}

// Create credentials from a filename
// Generally, prefer using this
func NewCredentialsFromFile(serverName string, filePath string) credentials.
	TransportCredentials {
	// Convert to fully qualified path
	filePath = utils.GetFullPath(filePath)
	// Generate credentials from path
	result, err := credentials.NewClientTLSFromFile(filePath, serverName)
	if err != nil {
		jww.FATAL.Panicf("Could not load TLS keys: %s", errors.New(err.Error()))
	}
	return result
}

type ConnectionManager struct {
	// A map of string IDs to open connections
	connections     map[string]*ConnectionInfo
	connectionsLock sync.Mutex
}

// Default maximum number of retries
const MAX_RETRIES = 5

// Convenience method to make a TransportCredentials for connecting
func MakeCreds(serverName, certPath, certPEM string) credentials.TransportCredentials {
	if certPath != "" {
		return NewCredentialsFromFile(serverName, certPath)
	} else if certPEM != "" {
		return NewCredentialsFromPEM(serverName, certPEM)
	} else {
		return nil
	}
}

// Connect to a certain registration server
// connectionInfo can be nil if the connection already exists for this id
func (m *ConnectionManager) ConnectToRegistration(id fmt.Stringer,
	info *ConnectionInfo) pb.RegistrationClient {
	connection := m.connect(id.String(), info)
	return pb.NewRegistrationClient(connection)
}

// Connect to a certain gateway
// connectionInfo can be nil if the connection already exists for this id
func (m *ConnectionManager) ConnectToGateway(id fmt.Stringer,
	info *ConnectionInfo) pb.GatewayClient {
	connection := m.connect(id.String(), info)
	return pb.NewGatewayClient(connection)
}

// Connect to a certain node
// connectionInfo can be nil if the connection already exists for this id
// Should this return an error if the connection doesn't exist and the
// connection info is nil?
func (m *ConnectionManager) ConnectToNode(id fmt.Stringer,
	info *ConnectionInfo) pb.NodeClient {
	connection := m.connect(id.String(), info)
	return pb.NewNodeClient(connection)
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

// Connect creates a connection, or returns a pre-existing connection based on
// a given address string.
// Connect should reconnect if the existing connection is non-nil,
// but the connection is no longer alive
// STILL UNDER CONSTRUCTION
func (m *ConnectionManager) connect(id string, info *ConnectionInfo) *grpc.
	ClientConn {
	// Check if a connection already exists
	m.connectionsLock.Lock() // TODO: Really we want to lock on the key,
	existingInfo, ok := m.connections[id]
	if ok && isConnectionGood(existingInfo.Connection) {
	}

	// Create top level vars
	var connection *grpc.ClientConn
	var err error
	connection = nil
	err = nil

	if m.connections == nil { // TODO: Do we need an init, or is this sufficient?
		m.connections = make(map[string]*ConnectionInfo)
	}

	maxRetries := 10
	// Create a new connection if we are not present or disconnecting/disconnected
	for numRetries := 0; numRetries < maxRetries && !isConnectionGood(connection); numRetries++ {

		jww.DEBUG.Printf("Trying to connect to %v", info.Address)
		ctx, cancel := context.WithTimeout(context.Background(),
			100000*time.Millisecond)

		if info.Creds != nil {
			// Create the GRPC client with TLS
			connection, err = grpc.DialContext(ctx, info.Address,
				grpc.WithTransportCredentials(info.Creds), grpc.WithBlock())
		} else {
			// Create the GRPC client without TLS
			connection, err = grpc.DialContext(ctx, info.Address,
				grpc.WithInsecure(), grpc.WithBlock())
		}

		if err == nil {
			// Connection succeeded; clean up context and exit the loop
			cancel()
		} else {
			jww.ERROR.Printf("Connection to %s failed: %+v\n", info.Address,
				errors.New(err.Error()))
		}
	}

	if !isConnectionGood(connection) {
		jww.FATAL.Panicf("Last try to connect to %s failed. Giving up",
			info.Address)
	} else {
		// Connection succeeded, so add it to the ConnectionInfo and record
		// the ConnectionInfo in the map
		info.Connection = connection
		m.connections[id] = info
	}

	m.connectionsLock.Unlock()

	return m.connections[id].Connection
}

// Disconnect closes client connections and removes them from the connection map
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

// DefaultContexts creates a context object with the default context
// for all client messages. This is primarily used to set the default
// timeout for all clients
func DefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(),
		10000*time.Millisecond)
	return ctx, cancel
}
