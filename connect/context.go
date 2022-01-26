///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Contains functionality related to context objects

package connect

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
	"net"
	"time"
)

// newContext builds a context for connections and message sending with the given timeout
func newContext(waitingPeriod time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(),
		waitingPeriod)
	return ctx, cancel
}

// StreamingContext creates a context object with the default context
// for all client streaming messages. This is primarily used to
// allow a cancel option for clients and is suitable for unary streaming.
func StreamingContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel
}

// StreamingContextWithTimeout creates a context object with the default
// context given a timout streaming messages. This is primarily used to
// allow a cancel option for clients and is suitable for unary streaming.
func StreamingContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return ctx, cancel
}

// GetAddressFromContext obtains address:port from the context of an incoming communication
func GetAddressFromContext(ctx context.Context) (address string, port string, err error) {
	info, ok := peer.FromContext(ctx)
	if !ok {
		return address, port, errors.New("Could not retrieve peer information from context")
	}
	address, port, err = net.SplitHostPort(info.Addr.String())
	return
}
