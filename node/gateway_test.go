////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package node

import (
	"context"
	"gitlab.com/elixxir/comms/connect"
	"gitlab.com/elixxir/comms/gateway"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"google.golang.org/grpc"
	"testing"
)

var GatewayAddress = "localhost:5557"

// Smoke test SendReceiveBatch
func TestSendReceiveBatch(t *testing.T) {
	gw := gateway.StartGateway(":5557",
		gateway.NewImplementation(), "", "")
	s := StartNode(":5558", NewImplementation(), "", "")
	defer gw.Shutdown()
	defer s.Shutdown()

	connID := MockID("mothership")
	s.ConnectToGateway(connID, &connect.ConnectionInfo{
		Address: ":5557",
	})
	err := s.SendReceiveBatch(connID, &pb.Batch{})
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
