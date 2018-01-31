// mixserver.go - Send/Receive functions for cMix servers
//
// Copyright © 2018 Privategrity Corporation
// All rights reserved.

package mixserver

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

// server object
type server struct {
	gs *grpc.Server
}

func ShutDown(s *server) {
	time.Sleep(time.Millisecond * 500)
	s.gs.GracefulStop()
}

// Handle a Broadcasted Network Error event
func (s *server) NetworkError(ctx context.Context, err *pb.ErrorMessage) (
	*pb.ErrorAck, error) {
	msgLen := int32(len(err.Message))
	jww.ERROR.Println(err.Message)
	return &pb.ErrorAck{MsgLen: msgLen}, nil
}

// Handle a Broadcasted Ask Online event
func (s *server) AskOnline(ctx context.Context, err *pb.Ping) (
	*pb.Pong, error) {
	return &pb.Pong{}, nil
}

// Send an AskOnline message to a particular server
func (s *server) SendAskOnline(addr string, message *pb.Ping) (*pb.Pong, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		jww.ERROR.Printf("Failed to connect to server with address %s",
			addr)
	}
	// Removed a sleep here

	c := pb.NewMixMessageServiceClient(conn)
	// Send AskOnline Request and check that we get an AskOnlineAck back
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	result, err = c.AskOnline(ctx, &pb.Ping{})
	if err != nil {
		jww.ERROR.Printf("AskOnline: Error received: %s", err)
	} else {
		jww.INFO.Printf("AskOnline: %v is online!", servers[i])
	}
	cancel()
	conn.Close()

	return result, nil
}

// Handle a PrecompDecrypt event
func (s *server) PrecompDecrypt(ctx context.Context, err *pb.PrecompDecryptMessage) (
	*pb.Ack, error) {
	return &pb.Ack{}, nil
}

func StartServer(addr string) {
	lis, err := net.Listen("tcp", addr)
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
