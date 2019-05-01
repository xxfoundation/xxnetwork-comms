////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
)

var GatewayAddress = "localhost:5557"

func startDummyGW() {
	// Listen on the given address
	lis, err := net.Listen("tcp", GatewayAddress)

	if err != nil {
		jww.FATAL.Panicf("failed to listen: %v", err)
	}

	//Make the port close when the gateway dies
	defer lis.Close()

	mixmessageServer := TestInterfaceGW{gs: grpc.NewServer()}
	pb.RegisterGatewayServer(mixmessageServer.gs, &mixmessageServer)

	// Register reflection service on gRPC server.
	// This blocks for the lifetime of the listener.
	reflection.Register(mixmessageServer.gs)
	if err := mixmessageServer.gs.Serve(lis); err != nil {
		jww.FATAL.Panicf("failed to serve: %v", err)
	}
}

// Smoke test SendReceiveBatch
func TestSendReceiveBatch(t *testing.T) {
	go startDummyGW()

	err := SendReceiveBatch(GatewayAddress, "", "", &pb.Batch{})
	if err != nil {
		t.Errorf("PutMessage: Error received: %s", err)
	}
}

// Blank struct implementing GatewayHandler interface for testing purposes
// (Passing to StartGateway)
type TestInterfaceGW struct {
	gs *grpc.Server
}

// Handle a CheckMessages event
func (s *TestInterfaceGW) CheckMessages(ctx context.Context, msg *pb.ClientRequest) (
	*pb.IDList, error) {
	returnMsg := &pb.IDList{}
	return returnMsg, nil
}

// Handle a GetMessage event
func (s *TestInterfaceGW) GetMessage(ctx context.Context, msg *pb.ClientRequest) (
	*pb.Batch, error) {
	returnMsg := &pb.Batch{}
	return returnMsg, nil
}

// Handle a PutMessage event
func (s *TestInterfaceGW) PutMessage(ctx context.Context, msg *pb.Slot) (*pb.Ack,
	error) {
	return &pb.Ack{}, nil
}

// Handle a PutMessage event
func (s *TestInterfaceGW) ReceiveBatch(ctx context.Context, msg *pb.Batch) (*pb.Ack,
	error) {
	return &pb.Ack{}, nil
}

// Pass-through for Registration Nonce Communication
func (s *TestInterfaceGW) RequestNonce(ctx context.Context,
	msg *pb.NonceRequest) (*pb.Nonce, error) {
	return &pb.Nonce{}, nil
}

// Pass-through for Registration Nonce Confirmation
func (s *TestInterfaceGW) ConfirmNonce(ctx context.Context,
	msg *pb.DSASignature) (*pb.RegistrationConfirmation, error) {
	return &pb.RegistrationConfirmation{}, nil
}
