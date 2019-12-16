////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server comms initialization functionality

package node

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

// Server object used to implement endpoints and top-level comms functionality
type Comms struct {
	connect.ProtoComms
	handler Handler
}

// Starts a new server on the address:port specified by listeningAddr
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNode(localServer string, handler Handler,
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
		// Create the TLS credentials
		// Create the TLS certificate
		x509cert, err2 := tls.X509KeyPair(certPEMblock, keyPEMblock)
		if err2 != nil {
			err = errors.New(err2.Error())
			jww.FATAL.Panicf("Could not load TLS keys: %+v", err)
		}

		creds := credentials.NewServerTLSFromCert(&x509cert)

		// Create the gRPC server with TLS
		jww.INFO.Printf("Starting server with TLS...")
		grpcServer = grpc.NewServer(grpc.Creds(creds),
			grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))
	} else {
		// Create the gRPC server without TLS
		jww.WARN.Printf("Starting server with TLS disabled...")
		grpcServer = grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
			grpc.MaxRecvMsgSize(math.MaxInt32))
	}

	mixmessageServer := Comms{
		ProtoComms: connect.ProtoComms{
			LocalServer:   grpcServer,
			ListeningAddr: localServer,
		},
		handler: handler,
	}

	go func() {
		// Register GRPC services to the listening address
		pb.RegisterNodeServer(mixmessageServer.LocalServer, &mixmessageServer)
		pb.RegisterGenericServer(mixmessageServer.LocalServer, &mixmessageServer)

		// Register reflection service on gRPC server.
		reflection.Register(mixmessageServer.LocalServer)
		if err = mixmessageServer.LocalServer.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
		jww.INFO.Printf("Shutting down node server listener: %s", lis)
	}()

	return &mixmessageServer
}
