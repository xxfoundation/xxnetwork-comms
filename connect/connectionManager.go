////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper/Helper functions for comms cMix client functionality
package connect

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

// A map of string addresses to open connections
var connections map[string]*grpc.ClientConn

// Holds the path for connecting to servers
// Must be explicitly set by gateways and servers to avoid data races
var ServerCertPath = ""

// A lock used to control access to the connections map above
var connectionsLock sync.Mutex

// Default maximum number of retries
const MAX_RETRIES = 5

// Connect to a gateway with a given address string
func ConnectToGateway(address string) pb.MixMessageGatewayClient {
	connection := connect(address, "")
	return pb.NewMixMessageGatewayClient(connection)
}

// Connect to a node with a given address string
func ConnectToNode(address string) pb.MixMessageNodeClient {
	connection := connect(address, ServerCertPath)
	return pb.NewMixMessageNodeClient(connection)
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
func connect(address string, certPath string) *grpc.ClientConn {
	var connection *grpc.ClientConn
	var err error
	connection = nil
	err = nil
	connectionsLock.Lock() // TODO: Really we want to lock on the key,
	// not the whole map

	if connections == nil { // TODO: Do we need an init, or is this sufficient?
		connections = make(map[string]*grpc.ClientConn)
	}

	maxRetries := 10
	// Create a new connection if we are not present or disconnecting/disconnected
	for numRetries := 0; numRetries < maxRetries && !isConnectionGood(
		address, connections); numRetries++ {

		jww.DEBUG.Printf("Trying to connect to %v", address)
		ctx, cancel := context.WithTimeout(context.Background(),
			100000*time.Millisecond)

		// If TLS was specified
		if certPath != "" {
			// Create the TLS credentials
			var creds credentials.TransportCredentials
			creds, err = credentials.NewClientTLSFromFile(certPath,
				"*.cmix.rip")
			if err != nil {
				jww.FATAL.Panicf("Could not load TLS keys: %s", err)
			}
			// Create the GRPC client with TLS
			connection, err = grpc.DialContext(ctx, address,
				grpc.WithTransportCredentials(creds), grpc.WithBlock())
		} else {
			// Create the GRPC client without TLS
			connection, err = grpc.DialContext(ctx, address,
				grpc.WithInsecure(), grpc.WithBlock())
		}

		if err == nil {
			connections[address] = connection
			cancel()
		} else {
			jww.ERROR.Printf("Connection to %s failed: %v\n", address, err)
		}
	}

	if !isConnectionGood(address, connections) {
		jww.FATAL.Panicf("Last try to connect to %s failed. Giving up", address)
	}

	connectionsLock.Unlock()

	return connections[address]
}

// Disconnect closes client connections and removes them from the connection map
func Disconnect(address string) {
	connectionsLock.Lock()
	connection, present := connections[address]
	if present {
		connection.Close()
		delete(connections, address)
	}
	connectionsLock.Unlock()
}

// DefaultContexts creates a context object with the default context
// for all client messages. This is primarily used to set the default
// timeout for all clients at 1/2 a second.
// TODO should gateway and node have different timeouts?
func DefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(),
		10000*time.Millisecond)
	return ctx, cancel
}
