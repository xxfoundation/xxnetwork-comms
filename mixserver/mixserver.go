// mixserver.go - Send/Receive functions for cMix servers
//
// Copyright © 2018 Privategrity Corporation
// All rights reserved.

package mixserver

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gitlab.com/privategrity/comms/mixmessages"
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	gs *grpc.Server
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (
	*pb.HelloReply, error) {
	// defer s.gs.GracefulStop()
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func StartServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	mixmessageServer := server{gs: grpc.NewServer()}
	pb.RegisterMixMessageServiceServer(mixmessageServer.gs, &mixmessageServer)
	// Register reflection service on gRPC server.
	reflection.Register(mixmessageServer.gs)
	if err := mixmessageServer.gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
