////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Contains registration server comms initialization functionality

package registration

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/comms/connect"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/comms/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
	"time"
)

var GlobalKeyPath string

// Server object containing a gRPC server
type RegistrationComms struct {
	connect.ConnectionManager
	gs      *grpc.Server
	handler Handler
}

// Performs a graceful shutdown of the server
func (s *RegistrationComms) Shutdown() {
	s.gs.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// Starts a new server on the address:port specified by localServer
// and a callback interface for server operations
// with given path to public and private key for TLS connection
func StartRegistrationServer(localServer string, handler Handler,
	certPath, keyPath string) *RegistrationComms {
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

		GlobalKeyPath = keyPath

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
	registrationServer := RegistrationComms{gs: grpcServer, handler: handler}

	go func() {
		// Make the port close when the gateway dies
		defer func() {
			err = lis.Close()
			if err != nil {
				err = errors.New(err.Error())
				jww.WARN.Printf("Unable to close listening port: %+v", err)
			}
		}()

		pb.RegisterRegistrationServer(registrationServer.gs, &registrationServer)

		// Register reflection service on gRPC server.
		reflection.Register(registrationServer.gs)
		if err = registrationServer.gs.Serve(lis); err != nil {
			err = errors.New(err.Error())
			jww.FATAL.Panicf("Failed to serve: %+v", err)
		}
	}()

	return &registrationServer
}
