////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	"google.golang.org/grpc"
	"net"
	"time"
)

// Passed into StartGateway to serve as an interface
// for interacting with the gateway repo
var gatewayHandler GatewayHandler

// gateway object
type gateway struct {
	gs *grpc.Server
}

func ShutDown(s *gateway) {
	time.Sleep(time.Millisecond * 500)
	s.gs.GracefulStop()
}

// Start local comm server
func StartGateway(localServer string, handler GatewayHandler) {
	// Set the gatewayHandler
	gatewayHandler = handler

	// Listen on the given address
	_, err := net.Listen("tcp", localServer)
	if err != nil {
		jww.FATAL.Panicf("failed to listen: %v", err)
	}
}
