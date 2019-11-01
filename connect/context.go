////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains functionality related to context objects

package connect

import (
	jww "github.com/spf13/jwalterweatherman"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

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

// Obtain address:port from the context of an incoming communication
func GetAddressFromContext(ctx context.Context) (address string, port string, err error) {
	info, _ := peer.FromContext(ctx)
	address, port, err = net.SplitHostPort(info.Addr.String())
	return
}
