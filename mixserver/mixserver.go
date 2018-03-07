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

// Handle a Broadcasted Network Error event
func (s *server) NetworkError(ctx context.Context, err *pb.ErrorMessage) (
	*pb.ErrorAck, error) {
	msgLen := int32(len(err.Message))
	jww.ERROR.Println(err.Message)
	return &pb.ErrorAck{MsgLen: msgLen}, nil
}

// Handle a Broadcasted Ask Online event
func (s *server) AskOnline(ctx context.Context, msg *pb.Ping) (
	*pb.Pong, error) {
	return &pb.Pong{}, nil
}

// Handle a NewRound event
func (s *server) NewRound(ctx context.Context,
	msg *pb.InitRound) (*pb.InitRoundAck, error) {
	// Call the server handler to start a new round
	serverHandler.NewRound(msg.RoundID)
	return &pb.InitRoundAck{}, nil
}

// Handle CmixMessage from Client to Server
func (s *server) ClientSendMessageToServer(ctx context.Context,
	msg *pb.CmixMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.ReceiveMessageFromClient(msg)
	return &pb.Ack{}, nil
}

// Request a CmixMessage from the server for the given User
func (s *server) ClientPoll(ctx context.Context,
	msg *pb.ClientPollMessage) (*pb.CmixMessage, error) {
	return serverHandler.ClientPoll(msg), nil
}

func (s *server) RequestContactList(ctx context.Context,
	msg *pb.ContactPoll) (*pb.ContactMessage, error) {
	return serverHandler.RequestContactList(msg), nil
}

// Handle a PrecompDecrypt event
func (s *server) PrecompDecrypt(ctx context.Context,
	msg *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompDecrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompEncrypt event
func (s *server) PrecompEncrypt(ctx context.Context,
	msg *pb.PrecompEncryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompEncrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompReveal event
func (s *server) PrecompReveal(ctx context.Context,
	msg *pb.PrecompRevealMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompReveal(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompPermute event
func (s *server) PrecompPermute(ctx context.Context,
	msg *pb.PrecompPermuteMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompPermute(msg)
	return &pb.Ack{}, nil
}

// Handle a PrecompShare event
func (s *server) PrecompShare(ctx context.Context,
	msg *pb.PrecompShareMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompShare(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimeDecrypt event
func (s *server) RealtimeDecrypt(ctx context.Context,
	msg *pb.RealtimeDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimeDecrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimeDecrypt event
func (s *server) RealtimeEncrypt(ctx context.Context,
	msg *pb.RealtimeEncryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimeEncrypt(msg)
	return &pb.Ack{}, nil
}

// Handle a RealtimePermute event
func (s *server) RealtimePermute(ctx context.Context,
	msg *pb.RealtimePermuteMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.RealtimePermute(msg)
	return &pb.Ack{}, nil
}

// Handle a SetPublicKey event
func (s *server) SetPublicKey(ctx context.Context,
	msg *pb.PublicKeyMessage) (*pb.Ack, error) {
	serverHandler.SetPublicKey(msg.RoundID, msg.PublicKey)
	return &pb.Ack{}, nil
}

// Starts the local comm server
func StartServer(localServer string, handler ServerHandler) {
	// Set the serverHandler
	serverHandler = handler

	// Listen on the given address
	lis, err := net.Listen("tcp", localServer)
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
