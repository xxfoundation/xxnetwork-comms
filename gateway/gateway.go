////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway comms initialization functionality

package gateway

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
	"time"
)

// Callback interface provided by the Gateway repository to StartGateway
var gatewayHandler Handler

// Gateway object containing a GRPC server
type gateway struct {
	gs *grpc.Server
}

// Performs a graceful shutdown of the gateway
func (s *gateway) ShutDown() {
	s.gs.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// Starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection
func StartGateway(localServer string, handler Handler,
	certPath, keyPath string) func() {
	var grpcServer *grpc.Server
	// Set the gatewayHandler
	gatewayHandler = handler

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		jww.FATAL.Panicf("Failed to listen: %v", err)
	}

	// If TLS was specified
	if certPath != "" && keyPath != "" {
		// Create the TLS credentials
		certPath = utils.GetFullPath(certPath)
		keyPath = utils.GetFullPath(keyPath)
		creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
		if err != nil {
			jww.FATAL.Panicf("Could not load TLS keys: %s", err)
		}

		// Create the GRPC server with TLS
		jww.INFO.Printf("Starting gateway with TLS...")
		grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(33554432)) // 32 MiB
	} else {
		// Create the GRPC server without TLS
		jww.INFO.Printf("Starting gateway with TLS disabled...")
		grpcServer = grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(33554432)) // 32 MiB
	}
	gatewayServer := gateway{gs: grpcServer}

	go func() {
		//Make the port close when the gateway dies
		defer func() {
			err := lis.Close()
			if err != nil {
				jww.WARN.Printf("Unable to close listening port: %s", err.Error())
			}
		}()

		pb.RegisterMixMessageGatewayServer(gatewayServer.gs, &gatewayServer)

		// Register reflection service on gRPC server.
		// This blocks for the lifetime of the listener.
		reflection.Register(gatewayServer.gs)
		if err := gatewayServer.gs.Serve(lis); err != nil {
			jww.FATAL.Panicf("Failed to serve: %v", err)
		}
	}()

	return gatewayServer.ShutDown
}
