////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains gateway comms initialization functionality

package gateway

import (
	"crypto/tls"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
)

// Gateway object used to implement endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms
	handler Handler
}

// Starts a new gateway on the address:port specified by localServer
// and a callback interface for gateway operations
// with given path to public and private key for TLS connection
func StartGateway(localServer string, handler Handler,
	certPEMblock, keyPEMblock []byte) *Comms {
	var grpcServer *grpc.Server

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		err = errors.New(err.Error())
		jww.FATAL.Panicf("Failed to listen: %+v", err)
	}

	// If TLS was specified
	if certPEMblock != nil && keyPEMblock != nil {
		// Create the TLS certificate
		x509cert, err2 := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err2 != nil {
			err = errors.New(err2.Error())
			jww.FATAL.Panicf("Could not load TLS keys: %+v", err)
		}

		creds := credentials.NewServerTLSFromCert(&x509cert)

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting gateway with TLS...")
		grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(33554432)) // 32 MiB

	} else {

		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting gateway with TLS disabled...")
		grpcServer = grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(33554432)) // 32 MiB

	}

	gatewayServer := Comms{
		ProtoComms: connect.ProtoComms{
			LocalServer:   grpcServer,
			ListeningAddr: localServer,
		},
		handler: handler,
	}

	go func() {
		pb.RegisterGatewayServer(gatewayServer.LocalServer, &gatewayServer)

		// Register reflection service on gRPC server.
		// This blocks for the lifetime of the listener.
		reflection.Register(gatewayServer.LocalServer)
		if err = gatewayServer.LocalServer.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
		jww.INFO.Printf("Shutting down gateway server listener: %s",
			lis)

	}()

	return &gatewayServer
}
