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

// Send an AskOnline message to a particular server
func SendAskOnline(addr string, message *pb.Ping) (*pb.Pong, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		jww.ERROR.Printf("Failed to connect to server with address %s",
			addr)
	}

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
func (s *server) PrecompDecrypt(ctx context.Context, msg *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Call the server handler with the msg
	serverHandler.PrecompDecrypt(msg)
	return &pb.Ack{}, nil
}

func SendPrecompDecrypt(nextServer string, input *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Create Dispatcher for Decrypt
	dcPrecompDecrypt := services.DispatchCryptop(node.Grp, precomputation.Decrypt{}, nil, nil, round)

	// Convert input message to equivalent SlotDecrypt
	slotDecrypt := &precomputation.SlotDecrypt{
		Slot:                         input.Slot,
		EncryptedMessageKeys:         cyclic.NewIntFromBytes(input.EncryptedMessageKeys),
		EncryptedRecipientIDKeys:     cyclic.NewIntFromBytes(input.EncryptedRecipientIDKeys),
		PartialMessageCypherText:     cyclic.NewIntFromBytes(input.PartialMessageCypherText),
		PartialRecipientIDCypherText: cyclic.NewIntFromBytes(input.PartialRecipientIDCypherText),
	}
	// Type assert SlotDecrypt to Slot
	var slot services.Slot = slotDecrypt

	// Pass slot as input to Decrypt
	dcPrecompDecrypt.InChannel <- &slot

	// Get output from Decrypt
	output := <-dcPrecompDecrypt.OutChannel
	// Type assert Slot to SlotDecrypt
	out := (*output).(*precomputation.SlotDecrypt)

	// Attempt to connect to nextServer
	conn, err := grpc.Dial(nextServer, grpc.WithInsecure())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", nextServer)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Send the PrecompDecrypt message using the Decrypt output
	result, err = c.PrecompDecrypt(ctx, &pb.PrecompDecryptMessage{
		Slot:                         out.Slot,
		EncryptedMessageKeys:         out.EncryptedMessageKeys.Bytes(),
		EncryptedRecipientIDKeys:     out.EncryptedRecipientIDKeys.Bytes(),
		PartialMessageCypherText:     out.PartialMessageCypherText.Bytes(),
		PartialRecipientIDCypherText: out.PartialRecipientIDCypherText.Bytes(),
	})
	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompDecrypt: Error received: %s", err)
	}

	return result, err
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
