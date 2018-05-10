package node

////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Wrapper/Helper functions for comms cMix client functionality

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"

	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/privategrity/comms/mixmessages"
)

// A map of string addresses to open connections
var connections map[string]*grpc.ClientConn

// A lock used to control access to the connections map above
var connectionsLock sync.Mutex

// Connect creates a connection, or returns a pre-existing connection based on
// a given address string.
func Connect(address string) pb.MixMessageServiceClient {
	var connection *grpc.ClientConn
	var err error
	connection = nil
	err = nil

	connectionsLock.Lock() // TODO: Really we want to lock on the key,
	// not the whole map

	if connections == nil { // TODO: Do we need an init, or is this sufficient?
		connections = make(map[string]*grpc.ClientConn)
	}

	// Check and return connection if it exists and is active
	connection, present := connections[address]

	// Create a new connection if we are not present or disconnecting/disconnected
	if !present || connection.GetState() == connectivity.Shutdown {
		for !present {
			ctx, cancel := context.WithTimeout(context.Background(),
				10000*time.Millisecond)
			connection, err = grpc.DialContext(ctx, address,
				grpc.WithInsecure(), grpc.WithBlock())
			if err == nil {
				connections[address] = connection
				cancel()
			} else {
				jww.WARN.Printf("Connection to %s failed, retrying: %v\n", address, err)
			}
			connection, present = connections[address]
		}
	}

	connectionsLock.Unlock()

	return pb.NewMixMessageServiceClient(connection)
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
func DefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	return ctx, cancel
}
