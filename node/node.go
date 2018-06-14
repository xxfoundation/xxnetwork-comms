////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package node handles all cMix node functionality. This file contains the
// main control logic when running a cMix Node.
package node

import (
	pb "gitlab.com/privategrity/comms/mixmessages"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	jww "github.com/spf13/jwalterweatherman"
	"net"
	"time"
	"math"
)

// Passed into StartServer to serve as an interface
// for interacting with the server repo
var serverHandler ServerHandler

// server object
type server struct {
	gs *grpc.Server
}

func (s *server) ShutDown() {
	s.gs.GracefulStop()
	time.Sleep(time.Millisecond * 500)
}

// StartServer starts a new server on the address:port specified by localServer
// NOTE: handler should be of type ServerImplementation. This will change
//       soon.
func StartServer(localServer string, handler ServerHandler) func() {
	// Set the serverHandler
	serverHandler = handler

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)

	if err != nil {
		jww.FATAL.Panicf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.MaxConcurrentStreams(math.MaxUint32),
		grpc.MaxRecvMsgSize(math.MaxInt64))

	mixmessageServer := server{gs: grpcServer}

	go func() {
		// Make the port close when the gateway dies
		defer func() {
			err := lis.Close()
			if err != nil {
				jww.WARN.Printf("Unable to close listening port: %s", err.Error())
			}
		}()

		pb.RegisterMixMessageNodeServer(mixmessageServer.gs, &mixmessageServer)

		// Register reflection service on gRPC server.
		reflection.Register(mixmessageServer.gs)
		if err := mixmessageServer.gs.Serve(lis); err != nil {
			jww.FATAL.Panicf("failed to serve: %v", err)
		}
	}()

	return mixmessageServer.ShutDown
}
