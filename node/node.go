////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains server comms initialization functionality

package node

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/utils"
	"math"
	"net"
	"time"
)

// Server object containing a gRPC server
type NodeComms struct {
	connect.ConnectionManager
	gs          *grpc.Server
	handler     ServerHandler
	localServer string
}

// Performs a graceful shutdown of the server
func (s *NodeComms) Shutdown() {
	// TODO Close all connections in the manager?
	s.gs.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartNode(localServer string, handler ServerHandler,
	certPath, keyPath, publicKey string) *NodeComms {
	var grpcServer *grpc.Server

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		err = errors.New(err.Error())
		jww.FATAL.Panicf("Failed to listen: %+v", err)
	}

	// If TLS was specified
	if certPath != "" && keyPath != "" {
		// Create the TLS credentials
		certPath = utils.GetFullPath(certPath)
		keyPath = utils.GetFullPath(keyPath)
		creds, err2 := credentials.NewServerTLSFromFile(certPath, keyPath)
		if err2 != nil {
			err = errors.New(err2.Error())
			jww.FATAL.Panicf("Could not load TLS keys: %+v", err)
		}

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
	mixmessageServer := NodeComms{gs: grpcServer, handler: handler, localServer: localServer}
	err = mixmessageServer.ConnectionManager.SetPublicKeyPath(publicKey)
	if err != nil {
		jww.ERROR.Printf("Error: %+v", err)
	}

	go func() {
		// Make the port close when the gateway dies
		defer func() {
			err = lis.Close()
			if err != nil {
				err = errors.New(err.Error())
				jww.WARN.Printf("Unable to close listening port: %+v", err)
			}
		}()

		pb.RegisterNodeServer(mixmessageServer.gs, &mixmessageServer)

		// Register reflection service on gRPC server.
		reflection.Register(mixmessageServer.gs)
		if err = mixmessageServer.gs.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
	}()

	return &mixmessageServer
}

func (s *NodeComms) String() string {
	return s.localServer
}
