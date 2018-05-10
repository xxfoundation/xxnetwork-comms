////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

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

func ShutDown(s *server) {
	time.Sleep(time.Millisecond * 500)
	s.gs.GracefulStop()
}

// Starts the local comm server
func StartServer(localServer string, handler ServerHandler) {
	// Set the serverHandler
	serverHandler = handler

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
	if err != nil {
		jww.FATAL.Panicf("failed to listen: %v", err)
	}
	mixmessageServer := server{gs: grpc.NewServer()}
	pb.RegisterMixMessageServiceServer(mixmessageServer.gs, &mixmessageServer)

	// Register reflection service on gRPC server.
	reflection.Register(mixmessageServer.gs)
	if err := mixmessageServer.gs.Serve(lis); err != nil {
		jww.FATAL.Panicf("failed to serve: %v", err)
	}
}
