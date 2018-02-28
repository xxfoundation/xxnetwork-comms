////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
//                                                                             /
////////////////////////////////////////////////////////////////////////////////
//                                                                             /
// Wrapper/Helper functions for comms cMix client functionality                /
//                                                                             /
////////////////////////////////////////////////////////////////////////////////
package mixclient

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"
	jww "github.com/spf13/jwalterweatherman"
)

// A map of string addresses to open connections
var connections map[string]*pb.MixMessageServiceClient
// A lock used to control access to the connections map above
var connectionsLock sync.Mutex

// Connect creates a connection, or returns a pre-existing connection based on
// a given address string.
func Connect(address string) (*pb.MixMessageServiceClient, error) {
	var connection *grpc.Clientconn
	var err *error
	connection = nil
	err = nil

	connectionsLock.Lock() // TODO: Really we want to lock on the key,
                         // not the whole map

	if connections == nil { // TODO: Do we need an init, or is this sufficient?
		connections = make(map[string]*grpc.ClientConn)
	}

	// Check and return connection if it exists and is active
	connection, present = connections[address]

	// Create a new connection if we are not present or disconnecting/disconnected
	if !present || connection.GetState() == connectivity.State.Shutdown {
		serverConn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
		if err == nil {
			connection = pb.NewMixMessageServiceClient(serverConn)
			connections[address] = connection
		} else {
			// TODO: Retry loop?
			jww.ERROR.Printf("Connection to %s failed: %v\n", address, err)
			connection = nil
		}
	}

	connectionsLock.Unlock()

	return connection
}

// DefaultContexts creates a context object with the default context
// for all client messages. This is primarily used to set the default
// timeout for all clients at 1/2 a second.
func DefaultContext() (context.Context, context.CancelFunc) {
	ctx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
	return ctx, cancel
}
