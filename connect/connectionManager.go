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
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
)

// A map of string addresses to open connections
var connections map[string]*grpc.ClientConn

// A lock used to control access to the connections map above
var connectionsLock sync.Mutex

// Default maximum number of retries
const MAX_RETRIES = 5

// Connect to a gateway with a given address string
func ConnectToGateway(address string) pb.MixMessageGatewayClient {
	connection := connect(address)
	return pb.NewMixMessageGatewayClient(connection)
}

// Connect to a node with a given address string
func ConnectToNode(address string) pb.MixMessageNodeClient {
	connection := connect(address)
	return pb.NewMixMessageNodeClient(connection)
}

// Is a connection in the map and alive?
func isConnectionGood(address string, connections map[string]*grpc.ClientConn) bool {
	connection, ok := connections[address]
	if !ok {
		return false
	}
	state := connection.GetState()
	return (state == connectivity.Idle || state == connectivity.Connecting ||
		state == connectivity.Ready)

}

// Connect creates a connection, or returns a pre-existing connection based on
// a given address string.
func connect(address string) *grpc.ClientConn {
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
		connection, err = grpc.DialContext(ctx, address,
			grpc.WithInsecure(), grpc.WithBlock())

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
