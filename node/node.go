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
)

// Passed into StartServer to serve as an interface
// for interacting with the server repo
var serverHandler ServerHandler

// server object
type server struct {
	gs *grpc.Server
}

var ServerObj *server

func ShutDown(s *server) {
	time.Sleep(time.Millisecond * 500)
	s.gs.GracefulStop()
}

// StartServer starts a new server on the address:port specified by localServer
// NOTE: handler should be of type ServerImplementation. This will change
//       soon.
func StartServer(localServer string, handler ServerHandler) {
	// Set the serverHandler
	serverHandler = handler

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)

	if err != nil {
		jww.FATAL.Panicf("failed to listen: %v", err)
	}

	// Make the port close when the gateway dies
	// This blocks for the lifetime of the listener.
	defer func() {
		err := lis.Close()
		if err != nil {
			jww.WARN.Printf("Unable to close listening port: %s", err.Error())
		}
	}()

	mixmessageServer := server{gs: grpc.NewServer()}
	ServerObj = &mixmessageServer
	pb.RegisterMixMessageNodeServer(mixmessageServer.gs, &mixmessageServer)

	// Register reflection service on gRPC server.
	reflection.Register(mixmessageServer.gs)
	if err := mixmessageServer.gs.Serve(lis); err != nil {
		jww.FATAL.Panicf("failed to serve: %v", err)
	}
}
