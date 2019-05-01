////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality for connecting to gateways and servers

package connect

import (
	"crypto/x509"
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
	Creds credentials.TransportCredentials
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
		jww.FATAL.Panicf("Could not load TLS keys: %s", err)
	}
	return result
}

type ConnectionManager struct {
	// A map of string addresses to open connections
	connections map[string]*grpc.ClientConn
	// TODO should this also store the credentials?
	connectionsLock sync.Mutex
}

// Default maximum number of retries
const MAX_RETRIES = 5

func makeCreds(serverName, certPath, certPEM string) credentials.
	TransportCredentials{
    if certPath != "" {
    	return NewCredentialsFromFile(serverName, certPath)
	} else if certPEM != "" {
		return NewCredentialsFromPEM(serverName, certPEM)
	} else {
		return nil
	}
}

// Connect to the registration server with a given address string
// TODO This should take a ConnectionInfo
func (m *ConnectionManager) ConnectToRegistration(address string,
	RegistrationCertPath string, RegistrationCertString string) pb.
	RegistrationClient {
	creds := makeCreds("registration*.cmix.rip", RegistrationCertPath, RegistrationCertString)
	connection := m.connect(&ConnectionInfo{
		Address: address,
		Creds:   creds,
	})
	return pb.NewRegistrationClient(connection)
}

// Connect to a gateway with a given address string
// TODO This should take a ConnectionInfo
func (m *ConnectionManager) ConnectToGateway(address string,
	GatewayCertPath string, GatewayCertString string) pb.GatewayClient {
	creds := makeCreds("gateway*.cmix.rip", GatewayCertPath, GatewayCertString)
	connection := m.connect(&ConnectionInfo{
		Address:address,
		Creds:creds,
	})
	return pb.NewGatewayClient(connection)
}

// Connect to a node with a given address string
// TODO This should take a ConnectionInfo
func (m *ConnectionManager) ConnectToNode(address string,
	ServerCertPath string) pb.NodeClient {
	creds := NewCredentialsFromFile("*.cmix.rip", ServerCertPath)
	connection := m.connect(&ConnectionInfo{
		Address: address,
		Creds:   creds,
	})
	return pb.NewNodeClient(connection)
}

// Is a connection in the map and alive?
func isConnectionGood(address string, connections map[string]*grpc.ClientConn) bool {
	connection, ok := connections[address]
	if !ok {
		return false
	}
	state := connection.GetState()
	return state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready

}

// Connect creates a connection, or returns a pre-existing connection based on
// a given address string.
func (m *ConnectionManager) connect(info *ConnectionInfo) *grpc.ClientConn {

	// Create top level vars
	var connection *grpc.ClientConn
	var err error
	connection = nil
	err = nil
	m.connectionsLock.Lock() // TODO: Really we want to lock on the key,

	if m.connections == nil { // TODO: Do we need an init, or is this sufficient?
		m.connections = make(map[string]*grpc.ClientConn)
	}

	maxRetries := 10
	// Create a new connection if we are not present or disconnecting/disconnected
	for numRetries := 0; numRetries < maxRetries && !isConnectionGood(
		info.Address, m.connections); numRetries++ {

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
			m.connections[info.Address] = connection
			cancel()
		} else {
			jww.ERROR.Printf("Connection to %s failed: %v\n", info.Address, err)
		}
	}

	if !isConnectionGood(info.Address, m.connections) {
		jww.FATAL.Panicf("Last try to connect to %s failed. Giving up",
			info.Address)
	}

	m.connectionsLock.Unlock()

	return m.connections[info.Address]
}

// Disconnect closes client connections and removes them from the connection map
func (m *ConnectionManager) Disconnect(address string) {
	m.connectionsLock.Lock()
	connection, present := m.connections[address]
	if present {
		err := connection.Close()
		if err != nil {
			jww.ERROR.Printf("Unable to close connection to %s: %v", address,
				err)
		}
		delete(m.connections, address)
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
